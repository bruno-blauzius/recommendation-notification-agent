package rabbitmq_test

import (
	"testing"

	"github.com/project-go-sender-recommendation-agent/internal/adapters/primary/rabbitmq"
)

func TestHelloWorldHandler_Handle_Success(t *testing.T) {
	h := rabbitmq.NewHelloWorldHandler()
	if err := h.Handle([]byte("some message")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHelloWorldHandler_Handle_EmptyPayload(t *testing.T) {
	h := rabbitmq.NewHelloWorldHandler()
	if err := h.Handle([]byte{}); err == nil {
		t.Fatal("expected error for empty payload")
	}
}

func TestHelloWorldHandler_Handle_LargePayload(t *testing.T) {
	h := rabbitmq.NewHelloWorldHandler()
	payload := make([]byte, 10*1024)
	for i := range payload {
		payload[i] = 'x'
	}
	if err := h.Handle(payload); err != nil {
		t.Fatalf("unexpected error for large payload: %v", err)
	}
}
