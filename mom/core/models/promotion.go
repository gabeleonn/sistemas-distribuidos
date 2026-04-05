package models

type Promotion struct {
	ID        string `json:"id"`
	Category  string `json:"category"`
	Item			string `json:"item"`
	HotDeal   bool   `json:"hot_deal"`
}

type PromotionReceivedPayload struct {
	Promotion Promotion `json:"promotion"`
}

type PromotionPublishedPayload struct {
	Promotion Promotion `json:"promotion"`
}

type PromotionVotePayload struct {
	Promotion Promotion `json:"promotion"`
	IsUpvote  bool      `json:"is_upvote"`
}

type PromotionFeaturedPayload struct {
	Promotion Promotion `json:"promotion"`
}
