package models

type Promotion struct {
	ID          string  `json:"id"`
	Store       string  `json:"store"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Score       float64 `json:"score"`
	Upvotes     int     `json:"upvotes"`
	Downvotes   int     `json:"downvotes"`
	IsHotDeal   bool    `json:"is_hot_deal"`
}

type PromotionReceivedPayload struct {
	Promotion Promotion `json:"promotion"`
}

type PromotionPublishedPayload struct {
	Promotion Promotion `json:"promotion"`
}

type PromotionVotePayload struct {
	PromotionID string `json:"promotion_id"`
	ClientID    int    `json:"client_id"`
	Positive    bool   `json:"positive"`
}

type PromotionFeaturedPayload struct {
	Promotion Promotion `json:"promotion"`
	Message   string    `json:"message"`
}
