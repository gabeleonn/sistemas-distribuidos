package notifications

import (
	"encoding/json"
	"fmt"
	"mom/core/amqp"
	"mom/core/constants"
	"mom/core/crypto"
	"mom/core/events"
	"mom/core/logger"
	"mom/core/models"

	"github.com/rabbitmq/amqp091-go"
)

func HandleConsumeCategories(client *amqp.Client) error {
	queue, err := client.DeclareQueue("", true, false ,false)
	if err != nil {
		return err
	}

	if err := client.BindQueue(queue.Name, constants.EventPromotionPublished); err != nil {
		return err
	}

	deliveries, err := client.Consume(queue.Name, constants.ServiceNotification)
	if err != nil {
		return err
	}

	for deliveries := range deliveries {
		var payload models.PromotionPublishedPayload
		if err := parseDelivery(deliveries, &payload); err != nil {
			continue
		}

		rawPayload, err := json.Marshal(payload.Promotion)

		if err != nil {
			continue
		}

		promokey := fmt.Sprintf("promocao.%s", payload.Promotion.Category)

		if err := client.Publish(promokey, rawPayload); err != nil {
			continue
		}

		logger.Get().Printf("promocao publicada: id=%s categoria=%s item=%s hotdeal=%t",
			payload.Promotion.ID,
			payload.Promotion.Category,
			payload.Promotion.Item,
			payload.Promotion.HotDeal,
		)
	}

	return nil

}

func parseDelivery(delivery amqp091.Delivery, payload *models.PromotionPublishedPayload) error {
	pkg, err := events.FromJSON(delivery.Body)
	if err != nil {
		return err
	}

	producerPubKey, err := crypto.LoadPublicKey(pkg.Producer)
	if err != nil {
		return err
	}

	if err :=  pkg.Verify(producerPubKey); err != nil {
		return err
	}

	return pkg.DecodePayload(payload)
}


