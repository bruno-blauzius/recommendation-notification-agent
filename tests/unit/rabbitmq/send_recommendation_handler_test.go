package rabbitmq_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/project-go-sender-recommendation-agent/internal/adapters/primary/rabbitmq"
	"github.com/project-go-sender-recommendation-agent/internal/core/domain"
)

// --- test doubles ---

type stubUseCase struct {
	createErr error
	createRec *domain.Recommendation
}

func (s *stubUseCase) Create(r *domain.Recommendation) error {
	s.createRec = r
	return s.createErr
}

func (s *stubUseCase) FindByID(_ string) (*domain.Recommendation, error) { return nil, nil }
func (s *stubUseCase) FindAll() ([]*domain.Recommendation, error)        { return nil, nil }

type stubIdempotency struct {
	existsResult bool
	existsErr    error
	setErr       error
	setCalled    bool
}

func (s *stubIdempotency) Exists(_ string) (bool, error) { return s.existsResult, s.existsErr }
func (s *stubIdempotency) Set(_ string, _ string, _ time.Duration) error {
	s.setCalled = true
	return s.setErr
}
func (s *stubIdempotency) Get(_ string) (string, error) { return "", nil }

func validPayload(t *testing.T, rec domain.Recommendation) []byte {
	t.Helper()
	b, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("failed to marshal recommendation: %v", err)
	}
	return b
}

// --- tests ---

func TestSendRecommendationHandler_Handle_Success(t *testing.T) {
	uc := &stubUseCase{}
	idm := &stubIdempotency{existsResult: false}
	h := rabbitmq.NewSendRecommendationHandler(uc, idm)

	payload := validPayload(t, domain.Recommendation{ID: "rec-1", SenderID: "s1", Payload: "data"})

	if err := h.Handle(payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uc.createRec == nil {
		t.Fatal("expected Create to be called")
	}
	if !idm.setCalled {
		t.Fatal("expected idempotency key to be set")
	}
}

func TestSendRecommendationHandler_Handle_EmptyPayload(t *testing.T) {
	h := rabbitmq.NewSendRecommendationHandler(&stubUseCase{}, &stubIdempotency{})
	if err := h.Handle([]byte{}); err == nil {
		t.Fatal("expected error for empty payload")
	}
}

func TestSendRecommendationHandler_Handle_InvalidJSON(t *testing.T) {
	h := rabbitmq.NewSendRecommendationHandler(&stubUseCase{}, &stubIdempotency{})
	if err := h.Handle([]byte("{invalid")); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSendRecommendationHandler_Handle_MissingID(t *testing.T) {
	h := rabbitmq.NewSendRecommendationHandler(&stubUseCase{}, &stubIdempotency{})
	payload := validPayload(t, domain.Recommendation{SenderID: "s1"}) // no ID
	if err := h.Handle(payload); err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestSendRecommendationHandler_Handle_DuplicateMessage_Skipped(t *testing.T) {
	uc := &stubUseCase{}
	idm := &stubIdempotency{existsResult: true} // already processing
	h := rabbitmq.NewSendRecommendationHandler(uc, idm)

	payload := validPayload(t, domain.Recommendation{ID: "rec-dup", SenderID: "s1"})

	if err := h.Handle(payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uc.createRec != nil {
		t.Fatal("expected Create to NOT be called for duplicate message")
	}
}

func TestSendRecommendationHandler_Handle_IdempotencyCheckError(t *testing.T) {
	idm := &stubIdempotency{existsErr: errors.New("redis error")}
	h := rabbitmq.NewSendRecommendationHandler(&stubUseCase{}, idm)

	payload := validPayload(t, domain.Recommendation{ID: "rec-1", SenderID: "s1"})
	if err := h.Handle(payload); err == nil {
		t.Fatal("expected error when idempotency check fails")
	}
}

func TestSendRecommendationHandler_Handle_IdempotencySetError(t *testing.T) {
	idm := &stubIdempotency{existsResult: false, setErr: errors.New("redis write error")}
	h := rabbitmq.NewSendRecommendationHandler(&stubUseCase{}, idm)

	payload := validPayload(t, domain.Recommendation{ID: "rec-1", SenderID: "s1"})
	if err := h.Handle(payload); err == nil {
		t.Fatal("expected error when idempotency set fails")
	}
}

func TestSendRecommendationHandler_Handle_UseCaseError(t *testing.T) {
	uc := &stubUseCase{createErr: errors.New("db error")}
	idm := &stubIdempotency{existsResult: false}
	h := rabbitmq.NewSendRecommendationHandler(uc, idm)

	payload := validPayload(t, domain.Recommendation{ID: "rec-1", SenderID: "s1"})
	if err := h.Handle(payload); err == nil {
		t.Fatal("expected error when use case fails")
	}
}
