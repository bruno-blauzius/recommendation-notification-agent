package rabbitmq

import (
	"fmt"
	"log"
)

// helloWorldHandler is a simple MessageHandler that logs the received payload.
type helloWorldHandler struct{}

// NewHelloWorldHandler returns a MessageHandler that prints "Hello World" for every message.
func NewHelloWorldHandler() *helloWorldHandler {
	return &helloWorldHandler{}
}

func (h *helloWorldHandler) Handle(payload []byte) error {
	log.Printf("Hello World — received message: %s", string(payload))
	if len(payload) == 0 {
		return fmt.Errorf("empty payload")
	}
	return nil
}
