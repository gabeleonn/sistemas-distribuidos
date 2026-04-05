package ranking

import (
	"mom/core/amqp"
	"mom/core/constants"
	"mom/core/crypto"
	"mom/core/events"
	"mom/core/logger"
	"mom/core/models"
	"mom/core/store"
	"os"
	"strconv"

	"github.com/rabbitmq/amqp091-go"
)

func ScoreHandler() error {
	logger.Init("ranking")

	storage := store.NewRankingStore()
	client := amqp.New()
	defer client.Close()

	queue, err := client.DeclareQueue("", true, false, false)
	if err != nil {
		return err
	}

	if err := client.BindQueue(queue.Name, constants.EventPromotionVote); err != nil {
		return err
	}

	deliveries, err := client.Consume(queue.Name, constants.ServiceRanking)
	if err != nil {
		return err
	}

	for delivery := range deliveries {
		var payload models.PromotionVotePayload
		if err := parseDelivery(delivery, &payload); err != nil {
			logger.Get().Errorf("erro ao decodificar payload do evento de voto de promocao: %v", err)
			continue
		}

		score := storage.ApplyVote(payload.Promotion.ID, payload.IsUpvote)
		
		threshold, err := strconv.Atoi(os.Getenv("HOTDEAL_SCORE_THRESHOLD"))
		if err != nil {
			logger.Get().Errorf("erro ao converter HOTDEAL_SCORE_THRESHOLD: %v", err)
			continue
		}

		if score == threshold {
			pkg, err := encodePromotionPublishedEvent(payload.Promotion)
			if err != nil {
				logger.Get().Errorf("erro ao codificar evento de promocao destacada: %v", err)
				continue
			}

			if err := client.Publish(constants.EventPromotionFeatured, pkg); err != nil {
				logger.Get().Errorf("erro ao publicar evento de promocao destacada: %v", err)
				continue
			}
			logger.Get().Printf("promocao %s atingiu pontuacao de hot deal: %d", payload.Promotion.ID, score)
		}
	}



	return nil
}

func encodePromotionPublishedEvent(promotion models.Promotion) ([]byte, error) {
	pkg, err := events.NewPackage(
		constants.ServiceRanking,
		constants.ServiceRanking,
		models.PromotionFeaturedPayload{
			Promotion: promotion,
		},
	)
	if err != nil {
		return nil, err
	}

	privateKey, _, err := crypto.EnsureKeyPair(constants.ServiceRanking)
	if err != nil {
		return nil, err
	}

	if err := pkg.Sign(privateKey); err != nil {
		return nil, err
	}

	return pkg.ToJSON()
}

func parseDelivery(delivery amqp091.Delivery, payload *models.PromotionVotePayload) error {
	pkg, err := events.FromJSON(delivery.Body)
	if err != nil {
		return err
	}

	producerPubKey, err := crypto.LoadPublicKey(pkg.Producer)
	if err != nil {
		return err
	}

	if err := pkg.Verify(producerPubKey); err != nil {
		return err
	}

	return pkg.DecodePayload(payload)
}
