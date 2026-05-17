package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/project-go-sender-recommendation-agent/internal/core/domain"
)

type recommendationRepository struct {
	db *sql.DB
}

// NewRecommendationRepository creates a new postgres-backed RecommendationRepository.
func NewRecommendationRepository(db *sql.DB) *recommendationRepository {
	return &recommendationRepository{db: db}
}

func (r *recommendationRepository) Save(rec *domain.Recommendation) error {
	payloadJSON, err := json.Marshal(rec.Payload)
	if err != nil {
		return fmt.Errorf("recommendationRepository.Save: marshal payload: %w", err)
	}

	query := `
		INSERT INTO recommendations (id, sender_id, payload, score, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	rec.CreatedAt = time.Now()
	_, err = r.db.Exec(query, rec.ID, rec.SenderID, payloadJSON, rec.Score, rec.CreatedAt)
	if err != nil {
		return fmt.Errorf("recommendationRepository.Save: %w", err)
	}
	return nil
}

func (r *recommendationRepository) FindByID(id string) (*domain.Recommendation, error) {
	query := `SELECT id, sender_id, payload, score, created_at FROM recommendations WHERE id = $1`
	row := r.db.QueryRow(query, id)

	rec := &domain.Recommendation{}
	var payloadJSON []byte
	err := row.Scan(&rec.ID, &rec.SenderID, &payloadJSON, &rec.Score, &rec.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("recommendationRepository.FindByID: %w", err)
	}
	if err := json.Unmarshal(payloadJSON, &rec.Payload); err != nil {
		return nil, fmt.Errorf("recommendationRepository.FindByID: unmarshal payload: %w", err)
	}
	return rec, nil
}

func (r *recommendationRepository) FindAll() ([]*domain.Recommendation, error) {
	query := `SELECT id, sender_id, payload, score, created_at FROM recommendations ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("recommendationRepository.FindAll: %w", err)
	}
	defer rows.Close()

	var recs []*domain.Recommendation
	for rows.Next() {
		rec := &domain.Recommendation{}
		var payloadJSON []byte
		if err := rows.Scan(&rec.ID, &rec.SenderID, &payloadJSON, &rec.Score, &rec.CreatedAt); err != nil {
			return nil, fmt.Errorf("recommendationRepository.FindAll scan: %w", err)
		}
		if err := json.Unmarshal(payloadJSON, &rec.Payload); err != nil {
			return nil, fmt.Errorf("recommendationRepository.FindAll: unmarshal payload: %w", err)
		}
		recs = append(recs, rec)
	}
	return recs, rows.Err()
}
