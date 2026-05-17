package rabbitmq

import (
	"crypto/md5"
	"encoding/hex"
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
	processingTTL       = 3 * time.Minute
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
	log.Printf("sendRecommendationHandler: received payload %s", string(payload))
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

	// MD5 of the raw payload guarantees uniqueness per message content
	// and avoids collisions that a bare ID alone could not prevent.
	hash := md5.Sum([]byte(payload))
	key := processingKeyPrefix + hex.EncodeToString(hash[:])

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
