package redis_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"

	redisadapter "github.com/project-go-sender-recommendation-agent/internal/adapters/secondary/redis"
)

func newTestRepo(t *testing.T) (*miniredis.Miniredis, interface {
	Set(string, string, time.Duration) error
	Get(string) (string, error)
	Exists(string) (bool, error)
}) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	t.Cleanup(mr.Close)

	client := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { client.Close() })

	return mr, redisadapter.NewIdempotencyRepository(client)
}

// --- Set ---

func TestIdempotencyRepository_Set_StoresValue(t *testing.T) {
	mr, repo := newTestRepo(t)

	if err := repo.Set("key:1", "processing", time.Minute); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := mr.Get("key:1")
	if err != nil {
		t.Fatalf("miniredis Get failed: %v", err)
	}
	if val != "processing" {
		t.Fatalf("expected 'processing', got %q", val)
	}
}

func TestIdempotencyRepository_Set_AppliesTTL(t *testing.T) {
	mr, repo := newTestRepo(t)

	if err := repo.Set("key:ttl", "1", 30*time.Second); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	ttl := mr.TTL("key:ttl")
	if ttl <= 0 {
		t.Fatalf("expected positive TTL, got %v", ttl)
	}
}

// --- Get ---

func TestIdempotencyRepository_Get_ReturnsValue(t *testing.T) {
	_, repo := newTestRepo(t)

	if err := repo.Set("key:get", "value-x", time.Minute); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := repo.Get("key:get")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "value-x" {
		t.Fatalf("expected 'value-x', got %q", val)
	}
}

func TestIdempotencyRepository_Get_KeyNotFound_ReturnsEmpty(t *testing.T) {
	_, repo := newTestRepo(t)

	val, err := repo.Get("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Fatalf("expected empty string, got %q", val)
	}
}

// --- Exists ---

func TestIdempotencyRepository_Exists_ReturnsTrueForExistingKey(t *testing.T) {
	_, repo := newTestRepo(t)

	if err := repo.Set("processing:rec-1", "1", time.Minute); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	exists, err := repo.Exists("processing:rec-1")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Fatal("expected key to exist")
	}
}

func TestIdempotencyRepository_Exists_ReturnsFalseForMissingKey(t *testing.T) {
	_, repo := newTestRepo(t)

	exists, err := repo.Exists("processing:missing")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Fatal("expected key to not exist")
	}
}

func TestIdempotencyRepository_Exists_ReturnsFalseAfterExpiry(t *testing.T) {
	mr, repo := newTestRepo(t)

	if err := repo.Set("processing:exp", "1", time.Second); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Fast-forward miniredis clock past TTL
	mr.FastForward(2 * time.Second)

	exists, err := repo.Exists("processing:exp")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Fatal("expected key to have expired")
	}
}
