package gateway

import (
	"mom/core/amqp"
	"mom/core/constants"
	"mom/core/crypto"
	"mom/core/events"
	"mom/core/logger"
	"mom/core/models"
	"mom/core/store"
)

func PublishedPromotionsHandler(storage *store.PromotionStore) error {
	client := amqp.New()
	defer client.Close()

	queue, err := client.DeclareQueue("", true, false, false)
	if err != nil {
		return err
	}

	if err := client.BindQueue(queue.Name, constants.EventPromotionPublished); err != nil {
		return err
	}

	deliveries, err := client.Consume(queue.Name, constants.ServiceGateway)
	if err != nil {
		return err
	}

	for delivery := range deliveries {
		pkg, err := events.FromJSON(delivery.Body)
		if err != nil {
			logger.Get().Errorf("erro ao ler evento publicado: %v", err)
			continue
		}

		publicKey, err := crypto.LoadPublicKey(pkg.Producer)
		if err != nil {
			logger.Get().Errorf("erro ao carregar chave publica: %v", err)
			continue
		}

		if err := pkg.Verify(publicKey); err != nil {
			logger.Get().Errorf("assinatura invalida: %v", err)
			continue
		}

		var payload models.PromotionPublishedPayload
		if err := pkg.DecodePayload(&payload); err != nil {
			logger.Get().Errorf("erro ao decodificar payload: %v", err)
			continue
		}

		storage.Save(payload.Promotion)
		logger.Get().Printf("promocao publicada recebida: id=%s categoria=%s", payload.Promotion.ID, payload.Promotion.Category)
	}

	return nil
}
