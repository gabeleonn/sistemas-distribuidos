package models

type Promotion struct {
	ID        string `json:"id"`
	Category  string `json:"category"`
	Item			string `json:"item"`
}

type PromotionReceivedPayload struct {
	Promotion Promotion `json:"promotion"`
}

type PromotionPublishedPayload struct {
	Promotion Promotion `json:"promotion"`
}

type PromotionVotePayload struct {
	PromotionID string `json:"promotion_id"`
	IsUpvote    bool   `json:"is_upvote"`
}

type PromotionFeaturedPayload struct {
	Promotion Promotion `json:"promotion"`
	Message   string    `json:"message"`
}
