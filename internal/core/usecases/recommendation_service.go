package usecases

import (
	"fmt"

	"github.com/project-go-sender-recommendation-agent/internal/core/domain"
	"github.com/project-go-sender-recommendation-agent/internal/core/ports/output"
)

type recommendationService struct {
	repo output.RecommendationRepository
}

// NewRecommendationService creates a new instance of RecommendationUseCase.
func NewRecommendationService(repo output.RecommendationRepository) *recommendationService {
	return &recommendationService{repo: repo}
}

func (s *recommendationService) Create(recommendation *domain.Recommendation) error {
	if recommendation.SenderID == "" {
		return fmt.Errorf("sender_id is required")
	}
	return s.repo.Save(recommendation)
}

func (s *recommendationService) FindByID(id string) (*domain.Recommendation, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return s.repo.FindByID(id)
}

func (s *recommendationService) FindAll() ([]*domain.Recommendation, error) {
	return s.repo.FindAll()
}
