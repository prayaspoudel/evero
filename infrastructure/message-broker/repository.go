// Package messagebroker provides message broker functionality for publishing
// and consuming messages from various message broker backends including RabbitMQ and NATS.
package messagebroker

import (
	"context"
	"time"
)

// MessageBroker interface defines the contract for message broker management
type MessageBroker interface {
	// Connect establishes connection to the message broker
	Connect(ctx context.Context) error

	// Disconnect closes the connection to the message broker
	Disconnect(ctx context.Context) error

	// Publish sends a message to the specified topic/queue
	Publish(ctx context.Context, topic string, message []byte, options *PublishOptions) error

	// PublishJSON sends a JSON-encoded message to the specified topic/queue
	PublishJSON(ctx context.Context, topic string, message interface{}, options *PublishOptions) error

	// Subscribe subscribes to messages from the specified topic/queue
	Subscribe(ctx context.Context, topic string, handler MessageHandler, options *SubscribeOptions) error

	// Unsubscribe unsubscribes from the specified topic/queue
	Unsubscribe(ctx context.Context, topic string) error

	// CreateTopic creates a new topic/queue (if supported by the broker)
	CreateTopic(ctx context.Context, topic string, options *TopicOptions) error

	// DeleteTopic deletes a topic/queue (if supported by the broker)
	DeleteTopic(ctx context.Context, topic string) error

	// ListTopics returns a list of available topics/queues
	ListTopics(ctx context.Context) ([]string, error)

	// Ping checks if the message broker is accessible
	Ping(ctx context.Context) error

	// Close closes the message broker and releases resources
	Close() error

	// PublishBatch publishes multiple messages in a batch
	PublishBatch(ctx context.Context, messages []BatchMessage, options *PublishOptions) error

	// GetStats returns broker statistics
	GetStats(ctx context.Context) (*BrokerStats, error)
}

// MessageHandler is a function type for handling incoming messages
type MessageHandler func(ctx context.Context, message *Message) error

// Message represents a message received from the broker
type Message struct {
	ID         string            `json:"id"`
	Topic      string            `json:"topic"`
	Data       []byte            `json:"data"`
	Headers    map[string]string `json:"headers"`
	Timestamp  time.Time         `json:"timestamp"`
	Retry      int               `json:"retry"`
	MaxRetries int               `json:"max_retries"`

	// Broker-specific fields
	OriginalMessage interface{} `json:"-"` // Store original message for acking
}

// BatchMessage represents a message for batch publishing
type BatchMessage struct {
	Topic   string            `json:"topic"`
	Data    []byte            `json:"data"`
	Headers map[string]string `json:"headers"`
}

// PublishOptions contains options for publishing messages
type PublishOptions struct {
	Headers     map[string]string `json:"headers"`
	Persistent  bool              `json:"persistent"`   // For durability
	Priority    int               `json:"priority"`     // Message priority (0-9)
	TTL         time.Duration     `json:"ttl"`          // Time to live
	Delay       time.Duration     `json:"delay"`        // Delay before delivery
	ContentType string            `json:"content_type"` // Content type
}

// SubscribeOptions contains options for subscribing to messages
type SubscribeOptions struct {
	QueueName     string        `json:"queue_name"`     // Queue name (for RabbitMQ)
	Durable       bool          `json:"durable"`        // Durable subscription
	AutoAck       bool          `json:"auto_ack"`       // Auto-acknowledge messages
	Exclusive     bool          `json:"exclusive"`      // Exclusive subscription
	MaxRetries    int           `json:"max_retries"`    // Maximum retry attempts
	RetryDelay    time.Duration `json:"retry_delay"`    // Delay between retries
	Concurrency   int           `json:"concurrency"`    // Number of concurrent handlers
	PrefetchCount int           `json:"prefetch_count"` // Number of messages to prefetch
}

// TopicOptions contains options for creating topics/queues
type TopicOptions struct {
	Durable    bool                   `json:"durable"`     // Durable topic/queue
	AutoDelete bool                   `json:"auto_delete"` // Auto-delete when unused
	Exclusive  bool                   `json:"exclusive"`   // Exclusive topic/queue
	Arguments  map[string]interface{} `json:"arguments"`   // Additional arguments
}

// BrokerStats contains broker statistics
type BrokerStats struct {
	ConnectedClients  int                    `json:"connected_clients"`
	MessagesPublished int64                  `json:"messages_published"`
	MessagesConsumed  int64                  `json:"messages_consumed"`
	Topics            []TopicStats           `json:"topics"`
	Custom            map[string]interface{} `json:"custom"` // Broker-specific stats
}

// TopicStats contains topic-specific statistics
type TopicStats struct {
	Name              string `json:"name"`
	MessagesPublished int64  `json:"messages_published"`
	MessagesConsumed  int64  `json:"messages_consumed"`
	Subscribers       int    `json:"subscribers"`
	PendingMessages   int64  `json:"pending_messages"`
}

// BrokerConfig holds configuration for message broker backends
type BrokerConfig struct {
	// RabbitMQ configuration
	RabbitMQURL      string `json:"rabbitmq_url"`
	RabbitMQExchange string `json:"rabbitmq_exchange"`
	RabbitMQVHost    string `json:"rabbitmq_vhost"`

	// NATS configuration
	NATSURL     string   `json:"nats_url"`
	NATSCluster string   `json:"nats_cluster"`
	NATSServers []string `json:"nats_servers"`

	// Kafka configuration
	KafkaURL              string   `json:"kafka_url"`
	KafkaBrokers          []string `json:"kafka_brokers"`
	KafkaConsumerGroup    string   `json:"kafka_consumer_group"`
	KafkaSASLMechanism    string   `json:"kafka_sasl_mechanism"`
	KafkaSecurityProtocol string   `json:"kafka_security_protocol"`

	// Connection settings
	MaxReconnects int           `json:"max_reconnects"`
	ReconnectWait time.Duration `json:"reconnect_wait"`
	Timeout       time.Duration `json:"timeout"`
	PingInterval  time.Duration `json:"ping_interval"`
	MaxPingsOut   int           `json:"max_pings_out"`

	// Authentication
	Username    string `json:"username"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`

	// TLS/SSL
	TLSEnabled    bool   `json:"tls_enabled"`
	TLSCertFile   string `json:"tls_cert_file"`
	TLSKeyFile    string `json:"tls_key_file"`
	TLSCAFile     string `json:"tls_ca_file"`
	TLSSkipVerify bool   `json:"tls_skip_verify"`
}
