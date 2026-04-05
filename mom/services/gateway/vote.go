package gateway

import (
	"fmt"

	amqp "mom/core/amqp"
	"mom/core/constants"
	"mom/core/crypto"
	"mom/core/events"
	"mom/core/logger"
	"mom/core/menu"
	"mom/core/models"
	"mom/core/store"
)

func VotePromotionHandler(storage *store.PromotionStore) error {
	promotions := storage.List()
	if len(promotions) == 0 {
		logger.Get().Println("nenhuma promocao publicada disponivel para voto")
		return nil
	}

	items := make([]string, 0, len(promotions))
	for _, promotion := range promotions {
		items = append(items, fmt.Sprintf("categoria=%s | item=%s | id=%s",
			promotion.Category,
			promotion.Item,
			promotion.ID,
		))	
	}

	index, _, err := menu.Select("Escolha a promocao para votar", items)
	if err != nil {
		return err
	}

	selectedPromotion := promotions[index]

	_, voteType, err := menu.Select("Escolha o tipo do voto", []string{
		"Positivo",
		"Negativo",
	})
	if err != nil {
		return err
	}

	isUpVote := voteType == "Positivo"
	vote := buildVote(selectedPromotion.ID, isUpVote)

	envelope, err := signVote(vote)
	if err != nil {
		return err
	}

	if err := publishVote(envelope); err != nil {
		return err
	}

	logger.Get().Printf("voto %q publicado para a promocao %q", voteType, selectedPromotion.ID)

	return nil
}

func buildVote(ID string, isUpVote bool) models.PromotionVotePayload {
	return models.PromotionVotePayload{
		PromotionID: ID,
		IsUpvote:    isUpVote,
	}
}

func signVote(vote models.PromotionVotePayload) (*events.Package, error) {
	producer := constants.ServiceGateway
	privateKey, _, err := crypto.EnsureKeyPair(producer)

	if err != nil {
		return nil, fmt.Errorf("erro ao carregar chave privada do consumidor: %w", err)
	}

	envelope, err := events.NewPackage(constants.EventPromotionVote, producer, vote)
	if err != nil {
		return nil, fmt.Errorf("erro ao montar evento de voto: %w", err)
	}

	if err := envelope.Sign(privateKey); err != nil {
		return nil, fmt.Errorf("erro ao assinar pacote de evento: %w", err)
	}

	return envelope, nil
}

func publishVote(envelope *events.Package) error {
	client := amqp.New()
	defer client.Close()

	body, err := envelope.ToJSON()
	if err != nil {
		return fmt.Errorf("erro ao converter pacote de voto para JSON: %w", err)
	}

	if err := client.Publish(constants.EventPromotionVote, body); err != nil {
		return fmt.Errorf("erro ao publicar evento de voto: %w", err)
	}

	return nil
}
