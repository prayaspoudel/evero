package messagebroker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// DefaultPublishOptions returns default publish options
func DefaultPublishOptions() *PublishOptions {
	return &PublishOptions{
		Headers:     make(map[string]string),
		Persistent:  false,
		Priority:    0,
		ContentType: "application/octet-stream",
	}
}

// DefaultSubscribeOptions returns default subscribe options
func DefaultSubscribeOptions() *SubscribeOptions {
	return &SubscribeOptions{
		AutoAck:       false,
		Durable:       true,
		Exclusive:     false,
		MaxRetries:    3,
		RetryDelay:    time.Second * 5,
		Concurrency:   1,
		PrefetchCount: 1,
	}
}

// DefaultTopicOptions returns default topic options
func DefaultTopicOptions() *TopicOptions {
	return &TopicOptions{
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		Arguments:  make(map[string]interface{}),
	}
}

// JSONPublishOptions returns publish options optimized for JSON messages
func JSONPublishOptions() *PublishOptions {
	opts := DefaultPublishOptions()
	opts.ContentType = "application/json"
	opts.Persistent = true
	return opts
}

// WorkerSubscribeOptions returns subscribe options optimized for worker queues
func WorkerSubscribeOptions(concurrency int) *SubscribeOptions {
	opts := DefaultSubscribeOptions()
	opts.Concurrency = concurrency
	opts.PrefetchCount = concurrency * 2
	opts.MaxRetries = 5
	opts.RetryDelay = time.Second * 10
	return opts
}

// MessageWithTimeout wraps a message handler with a timeout
func MessageWithTimeout(handler MessageHandler, timeout time.Duration) MessageHandler {
	return func(ctx context.Context, message *Message) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- handler(ctx, message)
		}()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return fmt.Errorf("message handler timeout after %v", timeout)
		}
	}
}

// MessageWithRetries wraps a message handler with custom retry logic
func MessageWithRetries(handler MessageHandler, maxRetries int, retryDelay time.Duration) MessageHandler {
	return func(ctx context.Context, message *Message) error {
		var lastErr error

		for attempt := 0; attempt <= maxRetries; attempt++ {
			err := handler(ctx, message)
			if err == nil {
				return nil
			}

			lastErr = err
			if attempt < maxRetries {
				select {
				case <-time.After(retryDelay):
					// Continue to next retry
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
	}
}

// JSONMessageHandler creates a handler that automatically unmarshals JSON messages
func JSONMessageHandler[T any](handler func(ctx context.Context, data T) error) MessageHandler {
	return func(ctx context.Context, message *Message) error {
		var data T
		if err := json.Unmarshal(message.Data, &data); err != nil {
			return fmt.Errorf("failed to unmarshal JSON message: %w", err)
		}

		return handler(ctx, data)
	}
}

// LoggingMessageHandler wraps a message handler with logging
func LoggingMessageHandler(handler MessageHandler, logger func(level string, msg string, args ...interface{})) MessageHandler {
	return func(ctx context.Context, message *Message) error {
		start := time.Now()
		logger("DEBUG", "Processing message", "topic", message.Topic, "id", message.ID)

		err := handler(ctx, message)
		duration := time.Since(start)

		if err != nil {
			logger("ERROR", "Message processing failed",
				"topic", message.Topic,
				"id", message.ID,
				"duration", duration.String(),
				"error", err.Error(),
			)
		} else {
			logger("INFO", "Message processed successfully",
				"topic", message.Topic,
				"id", message.ID,
				"duration", duration.String(),
			)
		}

		return err
	}
}

// BulkMessageHandler creates a handler that processes messages in batches
func BulkMessageHandler(batchSize int, flushInterval time.Duration, handler func(ctx context.Context, messages []*Message) error) MessageHandler {
	batch := make([]*Message, 0, batchSize)
	lastFlush := time.Now()
	var mutex sync.Mutex

	flushBatch := func(ctx context.Context) error {
		if len(batch) == 0 {
			return nil
		}

		err := handler(ctx, batch)
		batch = batch[:0] // Clear the batch
		lastFlush = time.Now()
		return err
	}

	return func(ctx context.Context, message *Message) error {
		mutex.Lock()
		defer mutex.Unlock()

		batch = append(batch, message)

		// Flush if batch is full or flush interval has passed
		shouldFlush := len(batch) >= batchSize || time.Since(lastFlush) >= flushInterval

		if shouldFlush {
			return flushBatch(ctx)
		}

		return nil
	}
}

// ValidateMessage validates a message structure
func ValidateMessage(message *Message) error {
	if message == nil {
		return errInvalidMessage
	}

	if message.Topic == "" {
		return fmt.Errorf("message topic cannot be empty")
	}

	if message.Data == nil {
		return fmt.Errorf("message data cannot be nil")
	}

	return nil
}

// CreateBatchMessage creates a batch message from individual components
func CreateBatchMessage(topic string, data interface{}, headers map[string]string) (BatchMessage, error) {
	var payload []byte
	var err error

	switch v := data.(type) {
	case []byte:
		payload = v
	case string:
		payload = []byte(v)
	default:
		payload, err = json.Marshal(data)
		if err != nil {
			return BatchMessage{}, fmt.Errorf("failed to marshal data: %w", err)
		}
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	return BatchMessage{
		Topic:   topic,
		Data:    payload,
		Headers: headers,
	}, nil
}

// MessageRouter provides simple message routing based on topic patterns
type MessageRouter struct {
	routes map[string]MessageHandler
}

// NewMessageRouter creates a new message router
func NewMessageRouter() *MessageRouter {
	return &MessageRouter{
		routes: make(map[string]MessageHandler),
	}
}

// AddRoute adds a route for a specific topic pattern
func (r *MessageRouter) AddRoute(pattern string, handler MessageHandler) {
	r.routes[pattern] = handler
}

// Route returns a message handler that routes messages based on topic
func (r *MessageRouter) Route() MessageHandler {
	return func(ctx context.Context, message *Message) error {
		// Simple exact match routing (could be extended to support patterns)
		if handler, exists := r.routes[message.Topic]; exists {
			return handler(ctx, message)
		}

		// Look for wildcard match
		if handler, exists := r.routes["*"]; exists {
			return handler(ctx, message)
		}

		return fmt.Errorf("no route found for topic: %s", message.Topic)
	}
}

// KafkaPublishOptions returns publish options optimized for Kafka messages
func KafkaPublishOptions(key string, partition int) *PublishOptions {
	opts := DefaultPublishOptions()
	opts.ContentType = "application/json"
	opts.Persistent = true

	if opts.Headers == nil {
		opts.Headers = make(map[string]string)
	}

	if key != "" {
		opts.Headers["kafka.key"] = key
	}

	if partition >= 0 {
		opts.Headers["kafka.partition"] = fmt.Sprintf("%d", partition)
	}

	return opts
}

// KafkaSubscribeOptions returns subscribe options optimized for Kafka consumers
func KafkaSubscribeOptions(consumerGroup string, concurrency int) *SubscribeOptions {
	opts := DefaultSubscribeOptions()
	opts.QueueName = consumerGroup // In Kafka, this becomes the consumer group
	opts.Concurrency = concurrency
	opts.PrefetchCount = concurrency * 2
	opts.MaxRetries = 3
	opts.RetryDelay = time.Second * 2
	opts.AutoAck = false // Kafka uses manual commit
	return opts
}

// KafkaTopicOptions returns topic options for creating Kafka topics
func KafkaTopicOptions(partitions, replicationFactor int) *TopicOptions {
	opts := DefaultTopicOptions()
	opts.Arguments["partitions"] = partitions
	opts.Arguments["replication-factor"] = replicationFactor
	return opts
}

// CircuitBreakerMessageHandler wraps a message handler with a simple circuit breaker
type CircuitBreakerMessageHandler struct {
	handler     MessageHandler
	maxFailures int
	resetTime   time.Duration
	failures    int
	lastFailure time.Time
	state       string // "closed", "open", "half-open"
}

// NewCircuitBreakerMessageHandler creates a new circuit breaker message handler
func NewCircuitBreakerMessageHandler(handler MessageHandler, maxFailures int, resetTime time.Duration) *CircuitBreakerMessageHandler {
	return &CircuitBreakerMessageHandler{
		handler:     handler,
		maxFailures: maxFailures,
		resetTime:   resetTime,
		state:       "closed",
	}
}

// Handle processes the message with circuit breaker logic
func (cb *CircuitBreakerMessageHandler) Handle(ctx context.Context, message *Message) error {
	switch cb.state {
	case "open":
		if time.Since(cb.lastFailure) > cb.resetTime {
			cb.state = "half-open"
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	}

	err := cb.handler(ctx, message)

	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()

		if cb.failures >= cb.maxFailures {
			cb.state = "open"
		}
	} else {
		if cb.state == "half-open" {
			cb.state = "closed"
			cb.failures = 0
		}
	}

	return err
}
