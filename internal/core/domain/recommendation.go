package domain

import "time"

// Recommendation represents the core entity of the recommendation domain.
type Recommendation struct {
	ID        string    `json:"id"`
	SenderID  string    `json:"sender_id"`
	Payload   string    `json:"payload"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}
