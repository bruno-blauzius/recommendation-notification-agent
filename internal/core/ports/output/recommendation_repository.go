package output

import "github.com/project-go-sender-recommendation-agent/internal/core/domain"

// RecommendationRepository defines the output port (driven port) for recommendation persistence.
type RecommendationRepository interface {
	Save(recommendation *domain.Recommendation) error
	FindByID(id string) (*domain.Recommendation, error)
	FindAll() ([]*domain.Recommendation, error)
}
