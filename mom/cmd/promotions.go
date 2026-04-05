package cmd

import (
	"fmt"
	"mom/core/amqp"
	"mom/core/constants"
	"mom/core/crypto"
	"mom/core/events"
	"mom/core/logger"
	"mom/core/models"
	"mom/core/store"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
)

var promoCmd = &cobra.Command{
	Use:   "promotions",
	Short: "Inicia o servico de promocoes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPromotions()
	},
}

func init() {
	rootCmd.AddCommand(promoCmd)
}

func runPromotions() error {
	log := logger.Init(constants.ServicePromotion)
	var storage = store.NewPromotionStore()

	client := amqp.New()
	defer client.Close()

	queue, err := client.DeclareQueue("", true, false, false)
	if err != nil {
		return fmt.Errorf("erro ao declarar fila de promocoes: %w", err)
	}

	if err := client.BindQueue(queue.Name, constants.EventPromotionReceived); err != nil {
		return fmt.Errorf("erro ao vincular fila de promocoes: %w", err)
	}

	deliveries, err := client.Consume(queue.Name, constants.ServicePromotion)

	if err != nil {
		return fmt.Errorf("erro ao consumir eventos de promocoes: %w", err)
	}

	for delivery := range deliveries {
		var payload models.PromotionReceivedPayload
		if err := parseDelivery(delivery, &payload); err != nil {
			log.Errorf("erro ao decodificar payload do evento de promocao: %v", err)
			continue
		}

		storage.Save(payload.Promotion)
		log.Printf("promocao recebida e salva: id=%s categoria=%s item=%s",
			payload.Promotion.ID,
			payload.Promotion.Category,
			payload.Promotion.Item,
		)

		body, err := encodePromotionPublishedEvent(payload.Promotion)
		if err != nil {
			log.Errorf("erro ao codificar evento de promocao publicada: %v", err)
			continue
		}

		if err := client.Publish(constants.EventPromotionPublished, body); err != nil {
			log.Errorf("erro ao publicar evento de promocao publicada: %v", err)
			continue
		}

		log.Printf("evento %q publicado para promocao %s",
			constants.EventPromotionPublished,
			payload.Promotion.ID,
		)
	}


	return nil
}

func parseDelivery(delivery amqp091.Delivery, payload *models.PromotionReceivedPayload) error {
	pkg, err := events.FromJSON(delivery.Body)
	if err != nil {
		return err
	}

	producerPubKey, err := crypto.LoadPublicKey(pkg.Producer)
	if err != nil {
		return err
	}
	
	if err := pkg.Verify(producerPubKey); err != nil {
		logger.Get().Errorf("erro ao verificar assinatura do evento de promocao: %v", err)
		return err
	}

	return pkg.DecodePayload(payload)
}

func encodePromotionPublishedEvent(promotion models.Promotion) ([]byte, error) {
	consumer := constants.ServicePromotion
	privateKey, _, err := crypto.EnsureKeyPair(consumer)

	if err != nil {
		return nil, err
	}

	newPkg, err := events.NewPackage(
		constants.EventPromotionPublished,
		consumer,
		models.PromotionPublishedPayload{
			Promotion: promotion,
		},
	)

	if err != nil {
		logger.Get().Errorf("erro ao criar pacote do evento de promocao publicada: %v", err)
		return nil, err
	}


	if err := newPkg.Sign(privateKey); err != nil {
		logger.Get().Errorf("erro ao assinar pacote do evento de promocao publicada: %v", err)
		return nil, err
	}

	body, err := newPkg.ToJSON()
	if err != nil {
		logger.Get().Errorf("erro ao converter pacote do evento de promocao publicada para JSON: %v", err)
		return nil, err
	}
	return body, nil
}
