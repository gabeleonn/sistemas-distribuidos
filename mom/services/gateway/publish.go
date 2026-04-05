package gateway

import (
	"fmt"
	amqpclient "mom/core/amqp"
	"mom/core/constants"
	"mom/core/crypto"
	"mom/core/events"
	"mom/core/logger"
	"mom/core/menu"
	"mom/core/models"

	"github.com/google/uuid"
)

func PublishPromotionHandler() error {
	categoria, err := menu.PromptText("Categoria da promocao")
	if err != nil {
		return err
	}

	item, err := menu.PromptText("Item da promocao")
	if err != nil {
		return err
	}

	promotion, err := buildPromotion(categoria, item)
	if err != nil {
		return err
	}

	envelope, err := signPromotion(promotion)
	if err != nil {
		return err
	}

	err = publishPromotion(envelope)
	if err != nil {
		return err
	}

	logger.Get().Printf("promocao com categoria %q publicada", categoria)
	return nil
}

func buildPromotion(category string, item string) (*models.Promotion, error) {
	return &models.Promotion{
		ID:        uuid.New().String(),
		Category:  category,
		Item: item,
	}, nil
}

func signPromotion(promotion *models.Promotion) (*events.Package, error) {
	producer := constants.ServiceGateway
	privateKey, _, err := crypto.EnsureKeyPair(producer)

	if err != nil {
		return nil, fmt.Errorf("erro ao carregar chave privada do consumidor: %w", err)
	}

	envelope, err := events.NewPackage(
		constants.EventPromotionReceived,
		producer,
		models.PromotionReceivedPayload{
			Promotion: *promotion,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar pacote de evento: %w", err)
	}

	if err := envelope.Sign(privateKey); err != nil {
		return nil, fmt.Errorf("erro ao assinar pacote de evento: %w", err)
	}

	return envelope, nil
}

func publishPromotion(promotion *events.Package) error {
	client := amqpclient.New()
	defer client.Close()

	body, err := promotion.ToJSON()
	if err != nil {
		return fmt.Errorf("erro ao serializar pacote de evento: %w", err)
	}

	if err := client.Publish(constants.EventPromotionReceived, body); err != nil {
		return fmt.Errorf("erro ao publicar evento de promocao recebida: %w", err)
	}

	return nil
}