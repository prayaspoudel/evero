package messagebroker

import (
	"fmt"
	"time"
)

// ConfigBuilder helps build BrokerConfig for different brokers
type ConfigBuilder struct {
	config *BrokerConfig
}

// NewConfigBuilder creates a new configuration builder
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &BrokerConfig{
			MaxReconnects: 3,
			ReconnectWait: 5 * time.Second,
			Timeout:       30 * time.Second,
			PingInterval:  30 * time.Second,
			MaxPingsOut:   2,
		},
	}
}

// WithAuth sets authentication credentials
func (cb *ConfigBuilder) WithAuth(username, password string) *ConfigBuilder {
	cb.config.Username = username
	cb.config.Password = password
	return cb
}

// WithTLS enables TLS/SSL configuration
func (cb *ConfigBuilder) WithTLS(enabled bool, certFile, keyFile, caFile string, skipVerify bool) *ConfigBuilder {
	cb.config.TLSEnabled = enabled
	cb.config.TLSCertFile = certFile
	cb.config.TLSKeyFile = keyFile
	cb.config.TLSCAFile = caFile
	cb.config.TLSSkipVerify = skipVerify
	return cb
}

// WithConnectionSettings sets connection-related settings
func (cb *ConfigBuilder) WithConnectionSettings(maxReconnects int, reconnectWait, timeout time.Duration) *ConfigBuilder {
	cb.config.MaxReconnects = maxReconnects
	cb.config.ReconnectWait = reconnectWait
	cb.config.Timeout = timeout
	return cb
}

// ForRabbitMQ configures the builder for RabbitMQ
func (cb *ConfigBuilder) ForRabbitMQ(url, exchange, vhost string) *ConfigBuilder {
	cb.config.RabbitMQURL = url
	cb.config.RabbitMQExchange = exchange
	cb.config.RabbitMQVHost = vhost
	return cb
}

// ForNATS configures the builder for NATS
func (cb *ConfigBuilder) ForNATS(url, cluster string, servers []string) *ConfigBuilder {
	cb.config.NATSURL = url
	cb.config.NATSCluster = cluster
	cb.config.NATSServers = servers
	return cb
}

// ForKafka configures the builder for Kafka
func (cb *ConfigBuilder) ForKafka(brokers []string, consumerGroup, saslMechanism string) *ConfigBuilder {
	cb.config.KafkaBrokers = brokers
	cb.config.KafkaConsumerGroup = consumerGroup
	cb.config.KafkaSASLMechanism = saslMechanism
	return cb
}

// WithKafkaURL sets a single Kafka broker URL (alternative to brokers list)
func (cb *ConfigBuilder) WithKafkaURL(url string) *ConfigBuilder {
	cb.config.KafkaURL = url
	return cb
}

// WithToken sets authentication token (for NATS)
func (cb *ConfigBuilder) WithToken(token string) *ConfigBuilder {
	cb.config.Token = token
	return cb
}

// Build returns the constructed BrokerConfig
func (cb *ConfigBuilder) Build() *BrokerConfig {
	return cb.config
}

// Predefined configuration builders for common setups

// LocalRabbitMQConfig creates a config for local RabbitMQ development
func LocalRabbitMQConfig(exchange string) *BrokerConfig {
	return NewConfigBuilder().
		ForRabbitMQ("amqp://localhost:5672", exchange, "/").
		Build()
}

// LocalNATSConfig creates a config for local NATS development
func LocalNATSConfig() *BrokerConfig {
	return NewConfigBuilder().
		ForNATS("nats://localhost:4222", "", nil).
		Build()
}

// LocalKafkaConfig creates a config for local Kafka development
func LocalKafkaConfig(consumerGroup string) *BrokerConfig {
	return NewConfigBuilder().
		ForKafka([]string{"localhost:9092"}, consumerGroup, "").
		Build()
}

// ProductionRabbitMQConfig creates a config for production RabbitMQ
func ProductionRabbitMQConfig(host, username, password, exchange string, port int, tlsEnabled bool) *BrokerConfig {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
	if tlsEnabled {
		url = fmt.Sprintf("amqps://%s:%s@%s:%d/", username, password, host, port)
	}

	builder := NewConfigBuilder().
		ForRabbitMQ(url, exchange, "/").
		WithAuth(username, password).
		WithConnectionSettings(5, 10*time.Second, 60*time.Second)

	if tlsEnabled {
		builder = builder.WithTLS(true, "", "", "", false)
	}

	return builder.Build()
}

// ProductionNATSConfig creates a config for production NATS
func ProductionNATSConfig(servers []string, username, password string, tlsEnabled bool) *BrokerConfig {
	builder := NewConfigBuilder().
		ForNATS("", "", servers).
		WithAuth(username, password).
		WithConnectionSettings(5, 10*time.Second, 60*time.Second)

	if tlsEnabled {
		builder = builder.WithTLS(true, "", "", "", false)
	}

	return builder.Build()
}

// ProductionKafkaConfig creates a config for production Kafka
func ProductionKafkaConfig(brokers []string, consumerGroup, username, password, saslMechanism string, tlsEnabled bool) *BrokerConfig {
	builder := NewConfigBuilder().
		ForKafka(brokers, consumerGroup, saslMechanism).
		WithAuth(username, password).
		WithConnectionSettings(5, 10*time.Second, 60*time.Second)

	if tlsEnabled {
		builder = builder.WithTLS(true, "", "", "", false)
	}

	return builder.Build()
}

// CloudKafkaConfig creates a config for cloud Kafka services (like Confluent Cloud)
func CloudKafkaConfig(brokers []string, consumerGroup, apiKey, apiSecret string) *BrokerConfig {
	return NewConfigBuilder().
		ForKafka(brokers, consumerGroup, "SASL_SSL").
		WithAuth(apiKey, apiSecret).
		WithTLS(true, "", "", "", false).
		WithConnectionSettings(5, 10*time.Second, 60*time.Second).
		Build()
}

// BrokerType represents the type of message broker
type BrokerType string

const (
	TypeRabbitMQ BrokerType = "rabbitmq"
	TypeNATS     BrokerType = "nats"
	TypeKafka    BrokerType = "kafka"
)

// CreateBroker is a convenience function to create a broker with config
func CreateBroker(brokerType BrokerType, config *BrokerConfig) (MessageBroker, error) {
	switch brokerType {
	case TypeRabbitMQ:
		return NewMessageBrokerFactory(InstanceRabbitMQ, config)
	case TypeNATS:
		return NewMessageBrokerFactory(InstanceNATS, config)
	case TypeKafka:
		return NewMessageBrokerFactory(InstanceKafka, config)
	default:
		return nil, fmt.Errorf("unsupported broker type: %s", brokerType)
	}
}

// AutoDetectBrokerType attempts to detect broker type from config
func AutoDetectBrokerType(config *BrokerConfig) BrokerType {
	if config.RabbitMQURL != "" {
		return TypeRabbitMQ
	}
	if config.NATSURL != "" || len(config.NATSServers) > 0 {
		return TypeNATS
	}
	if config.KafkaURL != "" || len(config.KafkaBrokers) > 0 {
		return TypeKafka
	}
	return ""
}

// CreateBrokerAuto automatically detects broker type and creates broker
func CreateBrokerAuto(config *BrokerConfig) (MessageBroker, error) {
	brokerType := AutoDetectBrokerType(config)
	if brokerType == "" {
		return nil, fmt.Errorf("could not detect broker type from config")
	}
	return CreateBroker(brokerType, config)
}
