package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/project-go-sender-recommendation-agent/internal/core/ports/input"
)

const (
	reconnectDelay = 5 * time.Second
	maxRetries     = 10
)

// Consumer is the primary adapter that reads messages from a RabbitMQ queue
// and delegates each message to a MessageHandler (input port).
type Consumer struct {
	dsn     string
	queue   string
	handler input.MessageHandler
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewConsumer dials RabbitMQ, declares the queue, and returns a ready Consumer.
func NewConsumer(dsn, queue string, handler input.MessageHandler) (*Consumer, error) {
	c := &Consumer{dsn: dsn, queue: queue, handler: handler}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Consumer) connect() error {
	conn, err := amqp.Dial(c.dsn)
	if err != nil {
		return fmt.Errorf("rabbitmq.Consumer dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("rabbitmq.Consumer channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		c.queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("rabbitmq.Consumer declare queue: %w", err)
	}

	c.conn = conn
	c.channel = ch
	return nil
}

// Consume starts consuming messages in a blocking loop with automatic reconnection.
// It retries up to maxRetries times on connection loss before returning an error.
func (c *Consumer) Consume() error {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		msgs, err := c.channel.Consume(
			c.queue,
			"",    // consumer tag (auto)
			false, // auto-ack disabled — ack only after successful handling
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("rabbitmq.Consumer.Consume register: %w", err)
		}

		log.Printf("RabbitMQ consumer started — queue: %q", c.queue)

		for msg := range msgs {
			if err := c.handler.Handle(msg.Body); err != nil {
				log.Printf("handler error: %v — nacking message", err)
				msg.Nack(false, true) // requeue on failure
				continue
			}
			msg.Ack(false)
		}

		// msgs channel closed: broker dropped the connection.
		log.Printf("RabbitMQ connection lost — reconnecting in %s (attempt %d/%d)", reconnectDelay, attempt, maxRetries)
		c.Close()
		time.Sleep(reconnectDelay)

		if err := c.connect(); err != nil {
			log.Printf("reconnect failed: %v", err)
		}
	}

	return fmt.Errorf("rabbitmq.Consumer.Consume: exhausted %d reconnect attempts", maxRetries)
}

// Close releases the channel and connection.
func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
