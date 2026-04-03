package cmd

import (
	"fmt"
	"mom/core/amqp"
	"mom/core/constants"
	"mom/core/crypto"
	"mom/core/events"
	"mom/core/models"
	"mom/core/store"

	"github.com/spf13/cobra"
)

var storage = store.NewPromotionStore()

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
	consumer := constants.ServicePromotion
	privateKey, _, err := crypto.EnsureKeyPair(consumer)

	if err != nil {
		return fmt.Errorf("erro ao carregar chave privada do consumidor: %w", err)
	}

	client := amqp.New()
	defer client.Close()

	queue, err := client.DeclareQueue("", true, false, false)
	if err != nil {
		return fmt.Errorf("erro ao declarar fila de promocoes: %w", err)
	}

	if err := client.BindQueue(queue.Name, constants.EventPromotionReceived); err != nil {
		return fmt.Errorf("erro ao vincular fila de promocoes: %w", err)
	}

	deliveries, err := client.Consume(queue.Name, consumer)

	if err != nil {
		return fmt.Errorf("erro ao consumir eventos de promocoes: %w", err)
	}

	for delivery := range deliveries {
		pkg, err := events.FromJSON(delivery.Body)
		if err != nil {
			fmt.Printf("erro ao decodificar evento de promocao: %v\n", err)
			continue
		}

		producerPubKey, err := crypto.LoadPublicKey(pkg.Producer)
		if err != nil {
			fmt.Printf("erro ao carregar chave publica do produtor: %v\n", err)
			continue
		}
		
		if err := pkg.Verify(producerPubKey); err != nil {
			fmt.Printf("erro ao verificar assinatura do evento de promocao: %v\n", err)
			continue
		}

		var payload models.PromotionReceivedPayload
		if err := pkg.DecodePayload(&payload); err != nil {
			fmt.Printf("erro ao decodificar payload do evento de promocao: %v\n", err)
			continue
		}

		storage.Save(payload.Promotion)

		newPkg, err := events.NewPackage(
			constants.EventPromotionPublished,
			consumer,
			models.PromotionPublishedPayload{
				Promotion: payload.Promotion,
			},
		)

		if err != nil {
			fmt.Printf("erro ao criar pacote do evento de promocao publicada: %v\n", err)
			continue
		}


		if err := newPkg.Sign(privateKey); err != nil {
			fmt.Printf("erro ao assinar pacote do evento de promocao publicada: %v\n", err)
			continue
		}

		body, err := newPkg.ToJSON()
		if err != nil {
			fmt.Printf("erro ao converter pacote do evento de promocao publicada para JSON: %v\n", err)
			continue
		}

		if err := client.Publish(constants.EventPromotionPublished, body); err != nil {
			fmt.Printf("erro ao publicar evento de promocao publicada: %v\n", err)
			continue
		}
	}


	return nil
}