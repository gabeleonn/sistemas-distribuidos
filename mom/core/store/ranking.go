package store

import "sync"

type RankingStore struct {
	mu     sync.RWMutex
	scores map[string]int
}

func NewRankingStore() *RankingStore {
	return &RankingStore{
		scores: make(map[string]int),
	}
}

func (s *RankingStore) ApplyVote(promotionID string, isUpvote bool) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if isUpvote {
		s.scores[promotionID]++
	} else {
		s.scores[promotionID]--
	}

	return s.scores[promotionID]
}

func (s *RankingStore) GetScore(promotionID string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	score, ok := s.scores[promotionID]
	return score, ok
}

func (s *RankingStore) Exists(promotionID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.scores[promotionID]
	return ok
}

func (s *RankingStore) Reset(promotionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.scores, promotionID)
}
