package gateway

import (
	"mom/core/logger"
	"mom/core/store"
)

func ShowPublishedPromotionsHandler(store *store.PromotionStore) error {
	promotions := store.List()
	if len(promotions) == 0 {
		logger.Get().Println("nenhuma promocao publicada por enquanto")
		return nil
	}

	logger.Get().Println("Promocoes publicadas:")
	for _, promo := range promotions {
		logger.Get().Printf("- ID: %s, Categoria: %s, Item: %s", promo.ID, promo.Category, promo.Item)
	}
	return nil
}