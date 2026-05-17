package usecases_test

import (
	"errors"
	"testing"

	"github.com/project-go-sender-recommendation-agent/internal/core/domain"
	"github.com/project-go-sender-recommendation-agent/internal/core/usecases"
)

// stubRecommendationRepo is a hand-written test double for output.RecommendationRepository.
type stubRecommendationRepo struct {
	saved      *domain.Recommendation
	findResult *domain.Recommendation
	findAllRes []*domain.Recommendation
	saveErr    error
	findErr    error
	findAllErr error
}

func (s *stubRecommendationRepo) Save(r *domain.Recommendation) error {
	s.saved = r
	return s.saveErr
}

func (s *stubRecommendationRepo) FindByID(_ string) (*domain.Recommendation, error) {
	return s.findResult, s.findErr
}

func (s *stubRecommendationRepo) FindAll() ([]*domain.Recommendation, error) {
	return s.findAllRes, s.findAllErr
}

// --- Create ---

func TestCreate_Success(t *testing.T) {
	stub := &stubRecommendationRepo{}
	svc := usecases.NewRecommendationService(stub)

	rec := &domain.Recommendation{ID: "1", SenderID: "sender-a", Payload: "data"}
	if err := svc.Create(rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.saved != rec {
		t.Fatal("expected Save to be called with the recommendation")
	}
}

func TestCreate_MissingSenderID(t *testing.T) {
	svc := usecases.NewRecommendationService(&stubRecommendationRepo{})
	err := svc.Create(&domain.Recommendation{ID: "1"})
	if err == nil {
		t.Fatal("expected error for empty sender_id")
	}
}

func TestCreate_RepoError(t *testing.T) {
	svc := usecases.NewRecommendationService(&stubRecommendationRepo{saveErr: errors.New("db error")})
	err := svc.Create(&domain.Recommendation{ID: "1", SenderID: "s1"})
	if err == nil {
		t.Fatal("expected repo error to propagate")
	}
}

// --- FindByID ---

func TestFindByID_Success(t *testing.T) {
	expected := &domain.Recommendation{ID: "1", SenderID: "s1"}
	svc := usecases.NewRecommendationService(&stubRecommendationRepo{findResult: expected})

	got, err := svc.FindByID("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != expected {
		t.Fatal("unexpected result")
	}
}

func TestFindByID_EmptyID(t *testing.T) {
	svc := usecases.NewRecommendationService(&stubRecommendationRepo{})
	if _, err := svc.FindByID(""); err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestFindByID_NotFound(t *testing.T) {
	svc := usecases.NewRecommendationService(&stubRecommendationRepo{findResult: nil})
	got, err := svc.FindByID("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Fatal("expected nil")
	}
}

func TestFindByID_RepoError(t *testing.T) {
	svc := usecases.NewRecommendationService(&stubRecommendationRepo{findErr: errors.New("db error")})
	if _, err := svc.FindByID("1"); err == nil {
		t.Fatal("expected repo error to propagate")
	}
}

// --- FindAll ---

func TestFindAll_Success(t *testing.T) {
	records := []*domain.Recommendation{{ID: "1"}, {ID: "2"}}
	svc := usecases.NewRecommendationService(&stubRecommendationRepo{findAllRes: records})

	got, err := svc.FindAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFindAll_Empty(t *testing.T) {
	svc := usecases.NewRecommendationService(&stubRecommendationRepo{})
	got, err := svc.FindAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFindAll_Error(t *testing.T) {
	svc := usecases.NewRecommendationService(&stubRecommendationRepo{findAllErr: errors.New("db error")})
	if _, err := svc.FindAll(); err == nil {
		t.Fatal("expected error from repo")
	}
}
