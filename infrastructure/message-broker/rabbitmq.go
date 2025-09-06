package messagebroker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQBroker struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	config      *BrokerConfig
	subscribers map[string]*rabbitMQSubscription
	mutex       sync.RWMutex
	connected   bool
}

type rabbitMQSubscription struct {
	channel  *amqp.Channel
	queue    string
	consumer string
	handler  MessageHandler
	options  *SubscribeOptions
	cancel   context.CancelFunc
	done     chan bool
}

// NewRabbitMQBroker creates a new RabbitMQ-based message broker
func NewRabbitMQBroker(config *BrokerConfig) (MessageBroker, error) {
	if config == nil {
		return nil, errors.New("broker config is required")
	}

	if config.RabbitMQURL == "" {
		return nil, errors.New("RabbitMQ URL is required")
	}

	return &rabbitMQBroker{
		config:      config,
		subscribers: make(map[string]*rabbitMQSubscription),
	}, nil
}

// Connect establishes connection to RabbitMQ
func (r *rabbitMQBroker) Connect(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var err error
	r.conn, err = amqp.Dial(r.config.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		r.conn.Close()
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	// Declare exchange if specified
	if r.config.RabbitMQExchange != "" {
		err = r.channel.ExchangeDeclare(
			r.config.RabbitMQExchange,
			"topic",
			true,  // durable
			false, // autoDelete
			false, // internal
			false, // noWait
			nil,   // arguments
		)
		if err != nil {
			r.cleanup()
			return fmt.Errorf("failed to declare exchange: %w", err)
		}
	}

	r.connected = true
	return nil
}

// Disconnect closes the RabbitMQ connection
func (r *rabbitMQBroker) Disconnect(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Stop all subscriptions
	for _, sub := range r.subscribers {
		if sub.cancel != nil {
			sub.cancel()
		}
		if sub.channel != nil {
			sub.channel.Close()
		}
	}
	r.subscribers = make(map[string]*rabbitMQSubscription)

	return r.cleanup()
}

func (r *rabbitMQBroker) cleanup() error {
	r.connected = false

	if r.channel != nil {
		r.channel.Close()
		r.channel = nil
	}

	if r.conn != nil {
		err := r.conn.Close()
		r.conn = nil
		return err
	}

	return nil
}

// Publish sends a message to the specified topic/queue
func (r *rabbitMQBroker) Publish(ctx context.Context, topic string, message []byte, options *PublishOptions) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if !r.connected {
		return errBrokerNotConnected
	}

	if options == nil {
		options = &PublishOptions{}
	}

	// Prepare publishing options
	publishing := amqp.Publishing{
		ContentType:  options.ContentType,
		Body:         message,
		Timestamp:    time.Now(),
		DeliveryMode: 1, // non-persistent
	}

	if options.Persistent {
		publishing.DeliveryMode = 2 // persistent
	}

	if options.Priority > 0 && options.Priority <= 9 {
		publishing.Priority = uint8(options.Priority)
	}

	if options.TTL > 0 {
		publishing.Expiration = fmt.Sprintf("%d", options.TTL.Milliseconds())
	}

	if options.Headers != nil {
		publishing.Headers = make(amqp.Table)
		for k, v := range options.Headers {
			publishing.Headers[k] = v
		}
	}

	// Publish to exchange or directly to queue
	var err error
	if r.config.RabbitMQExchange != "" {
		err = r.channel.Publish(
			r.config.RabbitMQExchange,
			topic, // routing key
			false, // mandatory
			false, // immediate
			publishing,
		)
	} else {
		// Ensure queue exists before publishing
		_, err = r.channel.QueueDeclare(
			topic,
			true,  // durable
			false, // autoDelete
			false, // exclusive
			false, // noWait
			nil,   // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue: %w", err)
		}

		err = r.channel.Publish(
			"",    // exchange
			topic, // routing key (queue name)
			false, // mandatory
			false, // immediate
			publishing,
		)
	}

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// PublishJSON sends a JSON-encoded message to the specified topic/queue
func (r *rabbitMQBroker) PublishJSON(ctx context.Context, topic string, message interface{}, options *PublishOptions) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if options == nil {
		options = &PublishOptions{}
	}
	if options.ContentType == "" {
		options.ContentType = "application/json"
	}

	return r.Publish(ctx, topic, data, options)
}

// Subscribe subscribes to messages from the specified topic/queue
func (r *rabbitMQBroker) Subscribe(ctx context.Context, topic string, handler MessageHandler, options *SubscribeOptions) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.connected {
		return errBrokerNotConnected
	}

	if options == nil {
		options = &SubscribeOptions{
			AutoAck:       false,
			MaxRetries:    3,
			RetryDelay:    time.Second,
			Concurrency:   1,
			PrefetchCount: 1,
		}
	}

	// Create a new channel for this subscription
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel for subscription: %w", err)
	}

	// Set QoS if prefetch count is specified
	if options.PrefetchCount > 0 {
		err = ch.Qos(options.PrefetchCount, 0, false)
		if err != nil {
			ch.Close()
			return fmt.Errorf("failed to set QoS: %w", err)
		}
	}

	queueName := topic
	if options.QueueName != "" {
		queueName = options.QueueName
	}

	// Declare queue
	queue, err := ch.QueueDeclare(
		queueName,
		options.Durable,   // durable
		false,             // autoDelete
		options.Exclusive, // exclusive
		false,             // noWait
		nil,               // arguments
	)
	if err != nil {
		ch.Close()
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange if exchange is configured
	if r.config.RabbitMQExchange != "" {
		err = ch.QueueBind(
			queue.Name,
			topic, // routing key
			r.config.RabbitMQExchange,
			false, // noWait
			nil,   // arguments
		)
		if err != nil {
			ch.Close()
			return fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	// Start consuming
	msgs, err := ch.Consume(
		queue.Name,
		"",                // consumer name (auto-generated)
		options.AutoAck,   // autoAck
		options.Exclusive, // exclusive
		false,             // noLocal
		false,             // noWait
		nil,               // arguments
	)
	if err != nil {
		ch.Close()
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	// Create subscription context
	subCtx, cancel := context.WithCancel(ctx)

	subscription := &rabbitMQSubscription{
		channel: ch,
		queue:   queue.Name,
		handler: handler,
		options: options,
		cancel:  cancel,
		done:    make(chan bool),
	}

	r.subscribers[topic] = subscription

	// Start message processing goroutines
	for i := 0; i < options.Concurrency; i++ {
		go r.processMessages(subCtx, msgs, subscription)
	}

	return nil
}

func (r *rabbitMQBroker) processMessages(ctx context.Context, msgs <-chan amqp.Delivery, subscription *rabbitMQSubscription) {
	defer func() {
		subscription.done <- true
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}

			r.handleMessage(ctx, msg, subscription)
		}
	}
}

func (r *rabbitMQBroker) handleMessage(ctx context.Context, delivery amqp.Delivery, subscription *rabbitMQSubscription) {
	message := &Message{
		ID:              delivery.MessageId,
		Topic:           delivery.RoutingKey,
		Data:            delivery.Body,
		Headers:         make(map[string]string),
		Timestamp:       delivery.Timestamp,
		OriginalMessage: delivery,
	}

	// Convert headers
	if delivery.Headers != nil {
		for k, v := range delivery.Headers {
			if str, ok := v.(string); ok {
				message.Headers[k] = str
			} else {
				message.Headers[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	// Process message with retries
	var lastErr error
	for retry := 0; retry <= subscription.options.MaxRetries; retry++ {
		message.Retry = retry
		message.MaxRetries = subscription.options.MaxRetries

		err := subscription.handler(ctx, message)
		if err == nil {
			// Success - acknowledge if not auto-ack
			if !subscription.options.AutoAck {
				delivery.Ack(false)
			}
			return
		}

		lastErr = err
		if retry < subscription.options.MaxRetries {
			time.Sleep(subscription.options.RetryDelay)
		}
	}

	// Failed after all retries
	if !subscription.options.AutoAck {
		delivery.Nack(false, false) // Don't requeue
	}

	// Log error (in a real implementation, you might want to send to a dead letter queue)
	fmt.Printf("Failed to process message after %d retries: %v\n", subscription.options.MaxRetries, lastErr)
}

// Unsubscribe unsubscribes from the specified topic/queue
func (r *rabbitMQBroker) Unsubscribe(ctx context.Context, topic string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	subscription, exists := r.subscribers[topic]
	if !exists {
		return errSubscriptionNotFound
	}

	if subscription.cancel != nil {
		subscription.cancel()
	}

	if subscription.channel != nil {
		subscription.channel.Close()
	}

	delete(r.subscribers, topic)
	return nil
}

// CreateTopic creates a new topic/queue
func (r *rabbitMQBroker) CreateTopic(ctx context.Context, topic string, options *TopicOptions) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if !r.connected {
		return errBrokerNotConnected
	}

	if options == nil {
		options = &TopicOptions{
			Durable: true,
		}
	}

	_, err := r.channel.QueueDeclare(
		topic,
		options.Durable,
		options.AutoDelete,
		options.Exclusive,
		false, // noWait
		amqp.Table(options.Arguments),
	)

	return err
}

// DeleteTopic deletes a topic/queue
func (r *rabbitMQBroker) DeleteTopic(ctx context.Context, topic string) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if !r.connected {
		return errBrokerNotConnected
	}

	_, err := r.channel.QueueDelete(topic, false, false, false)
	return err
}

// ListTopics returns a list of available topics/queues (limited in RabbitMQ)
func (r *rabbitMQBroker) ListTopics(ctx context.Context) ([]string, error) {
	// RabbitMQ doesn't provide a direct way to list all queues via AMQP
	// This would typically require the management API
	return nil, errors.New("listing topics not supported via AMQP protocol")
}

// Ping checks if RabbitMQ is accessible
func (r *rabbitMQBroker) Ping(ctx context.Context) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if !r.connected || r.conn == nil || r.conn.IsClosed() {
		return errBrokerNotConnected
	}

	return nil
}

// PublishBatch publishes multiple messages in a batch
func (r *rabbitMQBroker) PublishBatch(ctx context.Context, messages []BatchMessage, options *PublishOptions) error {
	for _, msg := range messages {
		err := r.Publish(ctx, msg.Topic, msg.Data, &PublishOptions{
			Headers:     msg.Headers,
			Persistent:  options.Persistent,
			Priority:    options.Priority,
			TTL:         options.TTL,
			ContentType: options.ContentType,
		})
		if err != nil {
			return fmt.Errorf("failed to publish batch message to topic %s: %w", msg.Topic, err)
		}
	}
	return nil
}

// GetStats returns broker statistics
func (r *rabbitMQBroker) GetStats(ctx context.Context) (*BrokerStats, error) {
	// Basic stats - in a real implementation, you might use RabbitMQ Management API
	return &BrokerStats{
		ConnectedClients: 1, // Current connection
		Custom: map[string]interface{}{
			"type":        "RabbitMQ",
			"exchange":    r.config.RabbitMQExchange,
			"subscribers": len(r.subscribers),
		},
	}, nil
}

// Close closes the RabbitMQ broker
func (r *rabbitMQBroker) Close() error {
	return r.Disconnect(context.Background())
}
