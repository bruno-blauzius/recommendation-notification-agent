package input

import "github.com/project-go-sender-recommendation-agent/internal/core/domain"

// RecommendationUseCase defines the input port (driving port) for recommendation operations.
type RecommendationUseCase interface {
	Create(recommendation *domain.Recommendation) error
	FindByID(id string) (*domain.Recommendation, error)
	FindAll() ([]*domain.Recommendation, error)
}
