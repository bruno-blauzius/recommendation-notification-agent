package rabbitmq_test

import (
	"testing"

	"github.com/project-go-sender-recommendation-agent/internal/adapters/primary/rabbitmq"
)

// stubHandler is a minimal MessageHandler for consumer constructor tests.
type stubHandler struct{}

func (s *stubHandler) Handle(_ []byte) error { return nil }

func TestNewConsumer_InvalidDSN_ReturnsError(t *testing.T) {
	_, err := rabbitmq.NewConsumer("amqp://invalid-host:5672/", "test-queue", &stubHandler{})
	if err == nil {
		t.Fatal("expected error for unreachable broker")
	}
}

func TestNewConsumer_EmptyDSN_ReturnsError(t *testing.T) {
	_, err := rabbitmq.NewConsumer("", "test-queue", &stubHandler{})
	if err == nil {
		t.Fatal("expected error for empty DSN")
	}
}
