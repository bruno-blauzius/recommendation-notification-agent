package integration_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"

	redisadapter "github.com/project-go-sender-recommendation-agent/internal/adapters/secondary/redis"
	"github.com/project-go-sender-recommendation-agent/internal/core/domain"
	"github.com/project-go-sender-recommendation-agent/internal/core/usecases"
)

// inMemoryRepo is a stateful test double that simulates a real persistence layer.
type inMemoryRepo struct {
	records map[string]*domain.Recommendation
}

func newInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{records: make(map[string]*domain.Recommendation)}
}

func (r *inMemoryRepo) Save(rec *domain.Recommendation) error {
	r.records[rec.ID] = rec
	return nil
}

func (r *inMemoryRepo) FindByID(id string) (*domain.Recommendation, error) {
	rec, ok := r.records[id]
	if !ok {
		return nil, nil
	}
	return rec, nil
}

func (r *inMemoryRepo) FindAll() ([]*domain.Recommendation, error) {
	all := make([]*domain.Recommendation, 0, len(r.records))
	for _, v := range r.records {
		all = append(all, v)
	}
	return all, nil
}

// TestCreateAndRetrieve_Integration exercises the full create → find-by-id → find-all flow.
func TestCreateAndRetrieve_Integration(t *testing.T) {
	repo := newInMemoryRepo()
	svc := usecases.NewRecommendationService(repo)

	rec := &domain.Recommendation{ID: "int-1", SenderID: "sender-a", Payload: "data", Score: 0.9}

	if err := svc.Create(rec); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := svc.FindByID("int-1")
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected to find created record")
	}
	if found.SenderID != "sender-a" {
		t.Fatalf("expected sender-a, got %s", found.SenderID)
	}

	all, err := svc.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 record, got %d", len(all))
	}
}

// TestIdempotencyFlow_Integration exercises the full idempotency guard flow using a real in-process Redis.
func TestIdempotencyFlow_Integration(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	defer client.Close()

	idempotencyRepo := redisadapter.NewIdempotencyRepository(client)
	repo := newInMemoryRepo()
	svc := usecases.NewRecommendationService(repo)

	const processingKey = "processing:rec-42"
	const processingTTL = 10 * time.Minute

	// First check: not processing yet.
	processing, err := idempotencyRepo.Exists(processingKey)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if processing {
		t.Fatal("expected record to not be in processing state initially")
	}

	// Mark as processing and create the recommendation.
	if err := idempotencyRepo.Set(processingKey, "1", processingTTL); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	rec := &domain.Recommendation{ID: "rec-42", SenderID: "sender-b", Payload: "payload"}
	if err := svc.Create(rec); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Duplicate delivery: idempotency check must prevent double processing.
	alreadyProcessing, err := idempotencyRepo.Exists(processingKey)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !alreadyProcessing {
		t.Fatal("expected processing key to still be set")
	}

	// Simulate TTL expiry: key must disappear after the window.
	mr.FastForward(processingTTL + time.Second)

	expired, err := idempotencyRepo.Exists(processingKey)
	if err != nil {
		t.Fatalf("Exists after expiry failed: %v", err)
	}
	if expired {
		t.Fatal("expected processing key to have expired after TTL")
	}
}

// TestCreateMultipleAndRetrieveAll_Integration verifies FindAll returns all saved records.
func TestCreateMultipleAndRetrieveAll_Integration(t *testing.T) {
	repo := newInMemoryRepo()
	svc := usecases.NewRecommendationService(repo)

	recs := []*domain.Recommendation{
		{ID: "a", SenderID: "s1", Payload: "p1"},
		{ID: "b", SenderID: "s2", Payload: "p2"},
		{ID: "c", SenderID: "s3", Payload: "p3"},
	}

	for _, r := range recs {
		if err := svc.Create(r); err != nil {
			t.Fatalf("Create failed for %s: %v", r.ID, err)
		}
	}

	all, err := svc.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 records, got %d", len(all))
	}
}
