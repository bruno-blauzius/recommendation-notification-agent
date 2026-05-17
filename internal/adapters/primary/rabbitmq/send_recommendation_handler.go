package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/project-go-sender-recommendation-agent/internal/core/domain"
	"github.com/project-go-sender-recommendation-agent/internal/core/ports/input"
	"github.com/project-go-sender-recommendation-agent/internal/core/ports/output"
)

const (
	processingKeyPrefix = "processing:"
	processingTTL       = 10 * time.Minute
)

type sendRecommendationHandler struct {
	useCase     input.RecommendationUseCase
	idempotency output.IdempotencyRepository
}

// NewSendRecommendationHandler returns a MessageHandler that processes
// messages from the send-recommendation queue with idempotency control.
func NewSendRecommendationHandler(useCase input.RecommendationUseCase, idempotency output.IdempotencyRepository) *sendRecommendationHandler {
	return &sendRecommendationHandler{useCase: useCase, idempotency: idempotency}
}

func (h *sendRecommendationHandler) Handle(payload []byte) error {
	if len(payload) == 0 {
		return fmt.Errorf("sendRecommendationHandler: empty payload")
	}

	var recommendation domain.Recommendation
	if err := json.Unmarshal(payload, &recommendation); err != nil {
		return fmt.Errorf("sendRecommendationHandler: unmarshal: %w", err)
	}

	if recommendation.ID == "" {
		return fmt.Errorf("sendRecommendationHandler: id is required")
	}

	// TODO: criar minha key usando o payload e convertendo em um MD5 ou SHA256 para evitar colisões e garantir unicidade
	key := processingKeyPrefix + recommendation.ID

	exists, err := h.idempotency.Exists(key)
	if err != nil {
		return fmt.Errorf("sendRecommendationHandler: idempotency check: %w", err)
	}
	if exists {
		log.Printf("sendRecommendationHandler: duplicate message skipped id=%s", recommendation.ID)
		return nil
	}

	if err := h.idempotency.Set(key, "1", processingTTL); err != nil {
		return fmt.Errorf("sendRecommendationHandler: set idempotency key: %w", err)
	}

	if err := h.useCase.Create(&recommendation); err != nil {
		return fmt.Errorf("sendRecommendationHandler: create: %w", err)
	}

	log.Printf("sendRecommendationHandler: processed id=%s sender_id=%s", recommendation.ID, recommendation.SenderID)
	return nil
}
