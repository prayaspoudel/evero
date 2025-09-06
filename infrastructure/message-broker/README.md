# Message Broker

This package provides a unified interface for message broker functionality with support for multiple backends including RabbitMQ, NATS, and Kafka.

## Features

- **Multiple Backends**: Support for RabbitMQ, NATS, and Kafka message brokers
- **Unified Interface**: Consistent API across all message broker backends
- **Context Support**: All operations support context for cancellation and timeouts
- **Publish/Subscribe**: Support for both publish/subscribe and queue-based messaging
- **Batch Operations**: Support for batch publishing
- **Retry Mechanism**: Built-in retry logic for message handling
- **Connection Management**: Proper connection handling with reconnection support
- **JSON Support**: Built-in JSON marshaling/unmarshaling
- **Headers Support**: Support for message headers and metadata

## Supported Backends

### RabbitMQ
- Full AMQP 0.9.1 protocol support
- Exchange and queue management
- Message persistence and TTL
- Dead letter queues support
- Connection pooling and retry mechanisms

### NATS
- High-performance messaging
- Subject-based messaging
- Queue groups for load balancing
- Lightweight and fast
- Built-in clustering support

### Kafka
- Distributed streaming platform
- Partitioned topics for scalability
- Consumer groups for load balancing
- Message ordering and durability
- High throughput and fault tolerance
- Stream processing capabilities

## Usage

### Basic Usage with Kafka

```go
// Create Kafka broker
config := &messagebroker.BrokerConfig{
    KafkaBrokers: []string{"localhost:9092"},
    KafkaConsumerGroup: "my-consumer-group",
}

broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceKafka, config)
if err != nil {
    log.Fatal(err)
}

// Connect to broker
ctx := context.Background()
err = broker.Connect(ctx)
if err != nil {
    log.Fatal(err)
}
defer broker.Close()

// Publish a message with key for partitioning
publishOptions := messagebroker.KafkaPublishOptions("user-123", 0)
err = broker.Publish(ctx, "user.events", []byte("Hello Kafka"), publishOptions)
if err != nil {
    log.Fatal(err)
}

// Subscribe to messages with consumer group
subscribeOptions := messagebroker.KafkaSubscribeOptions("my-consumer-group", 3)
handler := func(ctx context.Context, msg *messagebroker.Message) error {
    fmt.Printf("Received from partition %s: %s\n", 
        msg.Headers["kafka.partition"], string(msg.Data))
    return nil
}

err = broker.Subscribe(ctx, "user.events", handler, subscribeOptions)
if err != nil {
    log.Fatal(err)
}
```

### Basic Usage with RabbitMQ

```go
import "your-project/infrastructure/message-broker"

// Create RabbitMQ broker
config := &messagebroker.BrokerConfig{
    RabbitMQURL:      "amqp://localhost:5672",
    RabbitMQExchange: "my-exchange",
}

broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceRabbitMQ, config)
if err != nil {
    log.Fatal(err)
}

// Connect to broker
ctx := context.Background()
err = broker.Connect(ctx)
if err != nil {
    log.Fatal(err)
}
defer broker.Close()

// Publish a message
err = broker.Publish(ctx, "my.topic", []byte("Hello World"), nil)
if err != nil {
    log.Fatal(err)
}

// Subscribe to messages
handler := func(ctx context.Context, msg *messagebroker.Message) error {
    fmt.Printf("Received: %s\n", string(msg.Data))
    return nil
}

err = broker.Subscribe(ctx, "my.topic", handler, &messagebroker.SubscribeOptions{
    QueueName: "my-queue",
    Durable:   true,
})
if err != nil {
    log.Fatal(err)
}
```

### Basic Usage with NATS

```go
// Create NATS broker
config := &messagebroker.BrokerConfig{
    NATSURL: "nats://localhost:4222",
}

broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceNATS, config)
if err != nil {
    log.Fatal(err)
}

// Connect and use similar to RabbitMQ
ctx := context.Background()
err = broker.Connect(ctx)
if err != nil {
    log.Fatal(err)
}
defer broker.Close()

// Publishing and subscribing work the same way
```

### JSON Messages

```go
// Publish JSON message
data := map[string]interface{}{
    "user_id": 123,
    "action":  "login",
    "timestamp": time.Now(),
}

err = broker.PublishJSON(ctx, "user.events", data, &messagebroker.PublishOptions{
    Persistent: true,
    Headers: map[string]string{
        "version": "1.0",
    },
})

// Handle JSON messages
handler := func(ctx context.Context, msg *messagebroker.Message) error {
    var event map[string]interface{}
    err := json.Unmarshal(msg.Data, &event)
    if err != nil {
        return err
    }
    
    fmt.Printf("Event: %+v\n", event)
    return nil
}
```

### Batch Publishing

```go
messages := []messagebroker.BatchMessage{
    {
        Topic: "topic1",
        Data:  []byte("message 1"),
        Headers: map[string]string{"type": "info"},
    },
    {
        Topic: "topic2",
        Data:  []byte("message 2"),
        Headers: map[string]string{"type": "warning"},
    },
}

err = broker.PublishBatch(ctx, messages, &messagebroker.PublishOptions{
    Persistent: true,
})
```

### Advanced Subscription Options

```go
options := &messagebroker.SubscribeOptions{
    QueueName:     "worker-queue",
    Durable:       true,
    AutoAck:       false,
    MaxRetries:    3,
    RetryDelay:    time.Second * 5,
    Concurrency:   5,  // Process 5 messages concurrently
    PrefetchCount: 10, // Prefetch 10 messages
}

handler := func(ctx context.Context, msg *messagebroker.Message) error {
    // Process message
    fmt.Printf("Processing: %s (attempt %d/%d)\n", 
        string(msg.Data), msg.Retry+1, msg.MaxRetries+1)
    
    // Simulate work
    time.Sleep(time.Millisecond * 100)
    
    return nil
}

err = broker.Subscribe(ctx, "work.queue", handler, options)
```

## Configuration

### Kafka Configuration

```go
type BrokerConfig struct {
    // Kafka settings
    KafkaURL             string   `json:"kafka_url"`              // Single broker URL
    KafkaBrokers         []string `json:"kafka_brokers"`          // List of broker URLs
    KafkaConsumerGroup   string   `json:"kafka_consumer_group"`   // Consumer group ID
    KafkaSASLMechanism   string   `json:"kafka_sasl_mechanism"`   // SASL mechanism (PLAIN, SCRAM-SHA-256, SCRAM-SHA-512)
    KafkaSecurityProtocol string  `json:"kafka_security_protocol"` // Security protocol
    
    // Connection settings
    MaxReconnects   int           `json:"max_reconnects"`   // Max reconnection attempts
    ReconnectWait   time.Duration `json:"reconnect_wait"`   // Wait between reconnects
    Timeout         time.Duration `json:"timeout"`          // Connection timeout
    
    // Authentication
    Username    string `json:"username"`    // SASL username
    Password    string `json:"password"`    // SASL password
    
    // TLS/SSL
    TLSEnabled    bool   `json:"tls_enabled"`     // Enable TLS
    TLSCertFile   string `json:"tls_cert_file"`   // Client certificate
    TLSKeyFile    string `json:"tls_key_file"`    // Client private key
    TLSCAFile     string `json:"tls_ca_file"`     // CA certificate
    TLSSkipVerify bool   `json:"tls_skip_verify"` // Skip certificate verification
}
```

### RabbitMQ Configuration

```go
type BrokerConfig struct {
    // RabbitMQ settings
    RabbitMQURL      string `json:"rabbitmq_url"`      // AMQP connection URL
    RabbitMQExchange string `json:"rabbitmq_exchange"` // Exchange name
    RabbitMQVHost    string `json:"rabbitmq_vhost"`    // Virtual host
    
    // Connection settings
    MaxReconnects   int           `json:"max_reconnects"`   // Max reconnection attempts
    ReconnectWait   time.Duration `json:"reconnect_wait"`   // Wait between reconnects
    Timeout         time.Duration `json:"timeout"`          // Connection timeout
    
    // Authentication
    Username    string `json:"username"`    // Username
    Password    string `json:"password"`    // Password
    
    // TLS/SSL
    TLSEnabled    bool   `json:"tls_enabled"`     // Enable TLS
    TLSCertFile   string `json:"tls_cert_file"`   // Client certificate
    TLSKeyFile    string `json:"tls_key_file"`    // Client private key
    TLSCAFile     string `json:"tls_ca_file"`     // CA certificate
    TLSSkipVerify bool   `json:"tls_skip_verify"` // Skip certificate verification
}
```

### NATS Configuration

```go
type BrokerConfig struct {
    // NATS settings
    NATSURL     string   `json:"nats_url"`     // NATS server URL
    NATSCluster string   `json:"nats_cluster"` // Cluster name
    NATSServers []string `json:"nats_servers"` // List of NATS servers
    
    // Connection settings
    MaxReconnects int           `json:"max_reconnects"` // Max reconnection attempts
    ReconnectWait time.Duration `json:"reconnect_wait"` // Wait between reconnects
    PingInterval  time.Duration `json:"ping_interval"`  // Ping interval
    MaxPingsOut   int           `json:"max_pings_out"`  // Max outstanding pings
    
    // Authentication
    Username string `json:"username"` // Username
    Password string `json:"password"` // Password
    Token    string `json:"token"`    // Authentication token
    
    // TLS settings (same as RabbitMQ)
}
```

## Message Structure

```go
type Message struct {
    ID          string            `json:"id"`          // Message ID
    Topic       string            `json:"topic"`       // Topic/subject
    Data        []byte            `json:"data"`        // Message payload
    Headers     map[string]string `json:"headers"`     // Message headers
    Timestamp   time.Time         `json:"timestamp"`   // Message timestamp
    Retry       int               `json:"retry"`       // Current retry attempt
    MaxRetries  int               `json:"max_retries"` // Maximum retries allowed
}
```

## Error Handling

The package defines several error types:

- `errInvalidBrokerInstance`: Invalid broker backend type
- `errBrokerNotConnected`: Broker not connected
- `errTopicNotFound`: Topic/queue not found
- `errSubscriptionNotFound`: Subscription not found
- `errInvalidMessage`: Invalid message format
- `errPublishFailed`: Failed to publish message
- `errSubscribeFailed`: Failed to subscribe to topic

## Dependencies

### RabbitMQ Backend
```
go get github.com/rabbitmq/amqp091-go
```

### NATS Backend
```
go get github.com/nats-io/nats.go
```

### Kafka Backend
```
go get github.com/IBM/sarama
```

## Best Practices

1. **Always use context**: Pass appropriate context for timeouts and cancellation
2. **Handle errors gracefully**: Implement proper error handling in message handlers
3. **Use appropriate durability**: Set durable queues for important messages
4. **Implement idempotency**: Make message handlers idempotent when possible
5. **Monitor performance**: Use GetStats() to monitor broker performance
6. **Graceful shutdown**: Always call Close() when shutting down
7. **Use batch operations**: For high throughput, use batch publishing
8. **Set appropriate retry policies**: Configure retries based on your use case

## Monitoring and Observability

```go
// Get broker statistics
stats, err := broker.GetStats(ctx)
if err != nil {
    log.Printf("Failed to get stats: %v", err)
} else {
    log.Printf("Connected clients: %d", stats.ConnectedClients)
    log.Printf("Messages published: %d", stats.MessagesPublished)
    log.Printf("Messages consumed: %d", stats.MessagesConsumed)
}

// Check broker health
err = broker.Ping(ctx)
if err != nil {
    log.Printf("Broker health check failed: %v", err)
}
```
