package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// RecommendationItem represents a single insurance product recommendation.
type RecommendationItem struct {
	Produto         string  `json:"produto"`
	Ramo            string  `json:"ramo"`
	LogoURL         string  `json:"logo_url"`
	Seguradora      string  `json:"seguradora"`
	ScoreRelevancia float64 `json:"score_relevancia"`
	Valor           string  `json:"valor"`
	Justificativa   string  `json:"justificativa"`
}

// RecommendationPayload holds the structured content of a recommendation message.
type RecommendationPayload struct {
	ClienteDescricao   string               `json:"cliente_descricao"`
	PerfilIdentificado string               `json:"perfil_identificado"`
	Recomendacoes      []RecommendationItem `json:"recomendacoes"`
}

// Recommendation represents the core entity of the recommendation domain.
type Recommendation struct {
	ID        string                `json:"id"`
	SenderID  string                `json:"sender_id"`
	Payload   RecommendationPayload `json:"payload"`
	Score     float64               `json:"score"`
	CreatedAt time.Time             `json:"created_at"`
}

var timeFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02",
}

// UnmarshalJSON parses Recommendation from JSON, accepting created_at in
// multiple formats (RFC3339 with/without timezone, date-only).
func (r *Recommendation) UnmarshalJSON(data []byte) error {
	type alias struct {
		ID        string          `json:"id"`
		SenderID  string          `json:"sender_id"`
		Payload   json.RawMessage `json:"payload"`
		Score     float64         `json:"score"`
		CreatedAt string          `json:"created_at"`
	}

	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	r.ID = a.ID
	r.SenderID = a.SenderID
	r.Score = a.Score

	if len(a.Payload) > 0 && string(a.Payload) != "null" {
		if err := json.Unmarshal(a.Payload, &r.Payload); err != nil {
			return fmt.Errorf("recommendation: cannot parse payload: %w", err)
		}
	}

	if a.CreatedAt == "" {
		r.CreatedAt = time.Time{}
		return nil
	}

	for _, format := range timeFormats {
		if t, err := time.Parse(format, a.CreatedAt); err == nil {
			r.CreatedAt = t
			return nil
		}
	}

	return fmt.Errorf("recommendation: cannot parse created_at %q", a.CreatedAt)
}
