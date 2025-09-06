package messagebroker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type natsBroker struct {
	conn        *nats.Conn
	config      *BrokerConfig
	subscribers map[string]*natsSubscription
	mutex       sync.RWMutex
	connected   bool
}

type natsSubscription struct {
	subscription *nats.Subscription
	handler      MessageHandler
	options      *SubscribeOptions
	cancel       context.CancelFunc
	topic        string
}

// NewNATSBroker creates a new NATS-based message broker
func NewNATSBroker(config *BrokerConfig) (MessageBroker, error) {
	if config == nil {
		return nil, errors.New("broker config is required")
	}

	if config.NATSURL == "" && len(config.NATSServers) == 0 {
		return nil, errors.New("NATS URL or servers list is required")
	}

	return &natsBroker{
		config:      config,
		subscribers: make(map[string]*natsSubscription),
	}, nil
}

// Connect establishes connection to NATS
func (n *natsBroker) Connect(ctx context.Context) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	var opts []nats.Option

	// Set connection options
	if n.config.MaxReconnects > 0 {
		opts = append(opts, nats.MaxReconnects(n.config.MaxReconnects))
	}

	if n.config.ReconnectWait > 0 {
		opts = append(opts, nats.ReconnectWait(n.config.ReconnectWait))
	}

	if n.config.Timeout > 0 {
		opts = append(opts, nats.Timeout(n.config.Timeout))
	}

	if n.config.PingInterval > 0 {
		opts = append(opts, nats.PingInterval(n.config.PingInterval))
	}

	if n.config.MaxPingsOut > 0 {
		opts = append(opts, nats.MaxPingsOutstanding(n.config.MaxPingsOut))
	}

	// Authentication
	if n.config.Username != "" && n.config.Password != "" {
		opts = append(opts, nats.UserInfo(n.config.Username, n.config.Password))
	}

	if n.config.Token != "" {
		opts = append(opts, nats.Token(n.config.Token))
	}

	// TLS
	if n.config.TLSEnabled {
		if n.config.TLSCertFile != "" && n.config.TLSKeyFile != "" {
			opts = append(opts, nats.ClientCert(n.config.TLSCertFile, n.config.TLSKeyFile))
		}
		if n.config.TLSCAFile != "" {
			opts = append(opts, nats.RootCAs(n.config.TLSCAFile))
		}
		if n.config.TLSSkipVerify {
			opts = append(opts, nats.Secure())
		}
	}

	var err error
	var url string

	// Determine connection URL
	if len(n.config.NATSServers) > 0 {
		opts = append(opts, nats.DontRandomize())
		n.conn, err = nats.Connect(strings.Join(n.config.NATSServers, ","), opts...)
	} else {
		url = n.config.NATSURL
		if url == "" {
			url = nats.DefaultURL
		}
		n.conn, err = nats.Connect(url, opts...)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	n.connected = true
	return nil
}

// Disconnect closes the NATS connection
func (n *natsBroker) Disconnect(ctx context.Context) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// Stop all subscriptions
	for _, sub := range n.subscribers {
		if sub.cancel != nil {
			sub.cancel()
		}
		if sub.subscription != nil {
			sub.subscription.Unsubscribe()
		}
	}
	n.subscribers = make(map[string]*natsSubscription)

	if n.conn != nil {
		n.conn.Close()
		n.conn = nil
	}

	n.connected = false
	return nil
}

// Publish sends a message to the specified topic/queue
func (n *natsBroker) Publish(ctx context.Context, topic string, message []byte, options *PublishOptions) error {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	if !n.connected {
		return errBrokerNotConnected
	}

	// Create NATS message
	msg := &nats.Msg{
		Subject: topic,
		Data:    message,
	}

	if options != nil && options.Headers != nil {
		msg.Header = make(nats.Header)
		for k, v := range options.Headers {
			msg.Header.Set(k, v)
		}
	}

	// For NATS, we don't have built-in persistence or TTL like RabbitMQ
	// These would need to be handled at the application level or with NATS Streaming
	return n.conn.PublishMsg(msg)
}

// PublishJSON sends a JSON-encoded message to the specified topic/queue
func (n *natsBroker) PublishJSON(ctx context.Context, topic string, message interface{}, options *PublishOptions) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if options == nil {
		options = &PublishOptions{}
	}
	if options.Headers == nil {
		options.Headers = make(map[string]string)
	}
	options.Headers["Content-Type"] = "application/json"

	return n.Publish(ctx, topic, data, options)
}

// Subscribe subscribes to messages from the specified topic/queue
func (n *natsBroker) Subscribe(ctx context.Context, topic string, handler MessageHandler, options *SubscribeOptions) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if !n.connected {
		return errBrokerNotConnected
	}

	if options == nil {
		options = &SubscribeOptions{
			AutoAck:     true,
			MaxRetries:  3,
			RetryDelay:  time.Second,
			Concurrency: 1,
		}
	}

	subCtx, cancel := context.WithCancel(ctx)

	// Create NATS subscription
	var sub *nats.Subscription
	var err error

	// NATS message handler
	msgHandler := func(msg *nats.Msg) {
		n.handleNATSMessage(subCtx, msg, handler, options)
	}

	// Subscribe based on options
	if options.QueueName != "" {
		// Queue subscription (load balancing)
		sub, err = n.conn.QueueSubscribe(topic, options.QueueName, msgHandler)
	} else {
		// Regular subscription
		sub, err = n.conn.Subscribe(topic, msgHandler)
	}

	if err != nil {
		cancel()
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	// Store subscription
	natsSubscription := &natsSubscription{
		subscription: sub,
		handler:      handler,
		options:      options,
		cancel:       cancel,
		topic:        topic,
	}

	n.subscribers[topic] = natsSubscription
	return nil
}

func (n *natsBroker) handleNATSMessage(ctx context.Context, natsMsg *nats.Msg, handler MessageHandler, options *SubscribeOptions) {
	message := &Message{
		ID:              fmt.Sprintf("%d", time.Now().UnixNano()), // NATS doesn't have message IDs
		Topic:           natsMsg.Subject,
		Data:            natsMsg.Data,
		Headers:         make(map[string]string),
		Timestamp:       time.Now(),
		OriginalMessage: natsMsg,
	}

	// Convert headers
	if natsMsg.Header != nil {
		for k, v := range natsMsg.Header {
			if len(v) > 0 {
				message.Headers[k] = v[0]
			}
		}
	}

	// Process message with retries
	var lastErr error
	for retry := 0; retry <= options.MaxRetries; retry++ {
		message.Retry = retry
		message.MaxRetries = options.MaxRetries

		err := handler(ctx, message)
		if err == nil {
			// Success - NATS doesn't require explicit acking for regular subscriptions
			return
		}

		lastErr = err
		if retry < options.MaxRetries {
			time.Sleep(options.RetryDelay)
		}
	}

	// Failed after all retries
	fmt.Printf("Failed to process NATS message after %d retries: %v\n", options.MaxRetries, lastErr)
}

// Unsubscribe unsubscribes from the specified topic/queue
func (n *natsBroker) Unsubscribe(ctx context.Context, topic string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	subscription, exists := n.subscribers[topic]
	if !exists {
		return errSubscriptionNotFound
	}

	if subscription.cancel != nil {
		subscription.cancel()
	}

	if subscription.subscription != nil {
		err := subscription.subscription.Unsubscribe()
		if err != nil {
			return fmt.Errorf("failed to unsubscribe from topic %s: %w", topic, err)
		}
	}

	delete(n.subscribers, topic)
	return nil
}

// CreateTopic creates a new topic/queue (NATS doesn't require explicit topic creation)
func (n *natsBroker) CreateTopic(ctx context.Context, topic string, options *TopicOptions) error {
	// NATS doesn't require explicit topic creation
	// Topics are created dynamically when first published to
	return nil
}

// DeleteTopic deletes a topic/queue (NATS doesn't support topic deletion)
func (n *natsBroker) DeleteTopic(ctx context.Context, topic string) error {
	// NATS doesn't support explicit topic deletion
	// Topics are automatically cleaned up when no longer used
	return nil
}

// ListTopics returns a list of available topics/queues
func (n *natsBroker) ListTopics(ctx context.Context) ([]string, error) {
	// NATS doesn't provide a direct way to list all subjects
	// This would require using NATS monitoring or keeping track manually
	return nil, errors.New("listing topics not supported in NATS")
}

// Ping checks if NATS is accessible
func (n *natsBroker) Ping(ctx context.Context) error {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	if !n.connected || n.conn == nil || !n.conn.IsConnected() {
		return errBrokerNotConnected
	}

	return nil
}

// PublishBatch publishes multiple messages in a batch
func (n *natsBroker) PublishBatch(ctx context.Context, messages []BatchMessage, options *PublishOptions) error {
	// NATS doesn't have built-in batch publishing, so we publish one by one
	for _, msg := range messages {
		publishOptions := &PublishOptions{}
		if options != nil {
			*publishOptions = *options
		}
		if publishOptions.Headers == nil {
			publishOptions.Headers = make(map[string]string)
		}

		// Merge message headers with global options
		for k, v := range msg.Headers {
			publishOptions.Headers[k] = v
		}

		err := n.Publish(ctx, msg.Topic, msg.Data, publishOptions)
		if err != nil {
			return fmt.Errorf("failed to publish batch message to topic %s: %w", msg.Topic, err)
		}
	}
	return nil
}

// GetStats returns broker statistics
func (n *natsBroker) GetStats(ctx context.Context) (*BrokerStats, error) {
	stats := &BrokerStats{
		ConnectedClients: 1, // Current connection
		Custom: map[string]interface{}{
			"type":        "NATS",
			"subscribers": len(n.subscribers),
		},
	}

	if n.conn != nil {
		natsStats := n.conn.Stats()
		stats.Custom["in_msgs"] = natsStats.InMsgs
		stats.Custom["out_msgs"] = natsStats.OutMsgs
		stats.Custom["in_bytes"] = natsStats.InBytes
		stats.Custom["out_bytes"] = natsStats.OutBytes
		stats.Custom["reconnects"] = natsStats.Reconnects
	}

	return stats, nil
}

// Close closes the NATS broker
func (n *natsBroker) Close() error {
	return n.Disconnect(context.Background())
}
