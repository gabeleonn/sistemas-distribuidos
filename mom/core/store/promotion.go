package store

import (
	"sync"

	"mom/core/models"
)

type PromotionStore struct {
	mu         sync.RWMutex
	promotions map[string]models.Promotion
}

func NewPromotionStore() *PromotionStore {
	return &PromotionStore{
		promotions: make(map[string]models.Promotion),
	}
}

func (s *PromotionStore) Save(promotion models.Promotion) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.promotions[promotion.ID] = promotion
}

func (s *PromotionStore) GetByID(id string) (models.Promotion, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	promotion, ok := s.promotions[id]
	return promotion, ok
}

func (s *PromotionStore) Exists(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.promotions[id]
	return ok
}

func (s *PromotionStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.promotions, id)
}

func (s *PromotionStore) List() []models.Promotion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	promotions := make([]models.Promotion, 0, len(s.promotions))
	for _, promotion := range s.promotions {
		promotions = append(promotions, promotion)
	}

	return promotions
}

func (s *PromotionStore) ListFeatured() []models.Promotion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	promotions := make([]models.Promotion, 0)
	for _, promotion := range s.promotions {
		if promotion.IsHotDeal {
			promotions = append(promotions, promotion)
		}
	}

	return promotions
}

func (s *PromotionStore) Update(id string, updateFn func(*models.Promotion) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	promotion, ok := s.promotions[id]
	if !ok {
		return ErrPromotionNotFound
	}

	if err := updateFn(&promotion); err != nil {
		return err
	}

	s.promotions[id] = promotion
	return nil
}
