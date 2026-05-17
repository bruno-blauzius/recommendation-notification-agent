package input

// MessageHandler defines the input port for processing incoming messages.
// Any primary adapter (RabbitMQ, SQS, Kafka, etc.) must call this interface.
type MessageHandler interface {
	Handle(payload []byte) error
}
