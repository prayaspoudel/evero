package messagebroker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type kafkaBroker struct {
	config        *BrokerConfig
	producer      sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
	consumerGroup sarama.ConsumerGroup
	subscribers   map[string]*kafkaSubscription
	mutex         sync.RWMutex
	connected     bool
	client        sarama.Client
}

type kafkaSubscription struct {
	consumerGroup sarama.ConsumerGroup
	handler       MessageHandler
	options       *SubscribeOptions
	cancel        context.CancelFunc
	topic         string
	groupID       string
}

// kafkaConsumerGroupHandler implements sarama.ConsumerGroupHandler
type kafkaConsumerGroupHandler struct {
	subscription *kafkaSubscription
	broker       *kafkaBroker
}

// NewKafkaBroker creates a new Kafka-based message broker using Sarama
func NewKafkaBroker(config *BrokerConfig) (MessageBroker, error) {
	if config == nil {
		return nil, errors.New("broker config is required")
	}

	if len(config.KafkaBrokers) == 0 && config.KafkaURL == "" {
		return nil, errors.New("Kafka brokers list or URL is required")
	}

	return &kafkaBroker{
		config:      config,
		subscribers: make(map[string]*kafkaSubscription),
	}, nil
}

// Connect establishes connection to Kafka
func (k *kafkaBroker) Connect(ctx context.Context) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	brokers := k.getBrokers()

	// Create Sarama configuration
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_8_0_0 // Use a stable version

	// Producer configuration
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	// Consumer configuration
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Group.Session.Timeout = 10 * time.Second
	saramaConfig.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Authentication
	if k.config.Username != "" && k.config.Password != "" {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.User = k.config.Username
		saramaConfig.Net.SASL.Password = k.config.Password

		// Set SASL mechanism
		switch k.config.KafkaSASLMechanism {
		case "SCRAM-SHA-256":
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		case "SCRAM-SHA-512":
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		case "PLAIN":
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		default:
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		}
	}

	// TLS configuration
	if k.config.TLSEnabled {
		saramaConfig.Net.TLS.Enable = true
		// Additional TLS configuration can be added here
	}

	// Create client
	client, err := sarama.NewClient(brokers, saramaConfig)
	if err != nil {
		return fmt.Errorf("failed to create Kafka client: %w", err)
	}
	k.client = client

	// Create sync producer
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	k.producer = producer

	// Create async producer
	asyncProducer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		producer.Close()
		client.Close()
		return fmt.Errorf("failed to create Kafka async producer: %w", err)
	}
	k.asyncProducer = asyncProducer

	k.connected = true
	return nil
}

// Disconnect closes all Kafka connections
func (k *kafkaBroker) Disconnect(ctx context.Context) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	// Stop all subscriptions
	for _, sub := range k.subscribers {
		if sub.cancel != nil {
			sub.cancel()
		}
		if sub.consumerGroup != nil {
			sub.consumerGroup.Close()
		}
	}
	k.subscribers = make(map[string]*kafkaSubscription)

	// Close producers
	if k.asyncProducer != nil {
		k.asyncProducer.Close()
		k.asyncProducer = nil
	}

	if k.producer != nil {
		k.producer.Close()
		k.producer = nil
	}

	// Close client
	if k.client != nil {
		k.client.Close()
		k.client = nil
	}

	k.connected = false
	return nil
}

func (k *kafkaBroker) getBrokers() []string {
	if len(k.config.KafkaBrokers) > 0 {
		return k.config.KafkaBrokers
	}

	if k.config.KafkaURL != "" {
		return []string{k.config.KafkaURL}
	}

	return []string{"localhost:9092"} // default
}

// Publish sends a message to the specified topic
func (k *kafkaBroker) Publish(ctx context.Context, topic string, message []byte, options *PublishOptions) error {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	if !k.connected {
		return errBrokerNotConnected
	}

	// Create producer message
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(message),
		Timestamp: time.Now(),
	}

	if options != nil {
		// Add headers
		if options.Headers != nil {
			headers := make([]sarama.RecordHeader, 0, len(options.Headers))
			for k, v := range options.Headers {
				headers = append(headers, sarama.RecordHeader{
					Key:   []byte(k),
					Value: []byte(v),
				})
			}
			msg.Headers = headers
		}

		// Set key for partitioning if provided in headers
		if key, exists := options.Headers["kafka.key"]; exists {
			msg.Key = sarama.StringEncoder(key)
		}

		// Set partition if provided in headers
		if partitionStr, exists := options.Headers["kafka.partition"]; exists {
			if partition, err := strconv.Atoi(partitionStr); err == nil {
				msg.Partition = int32(partition)
			}
		}
	}

	// Send message synchronously
	_, _, err := k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	return nil
}

// PublishJSON sends a JSON-encoded message to the specified topic
func (k *kafkaBroker) PublishJSON(ctx context.Context, topic string, message interface{}, options *PublishOptions) error {
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

	return k.Publish(ctx, topic, data, options)
}

// Subscribe subscribes to messages from the specified topic
func (k *kafkaBroker) Subscribe(ctx context.Context, topic string, handler MessageHandler, options *SubscribeOptions) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if !k.connected {
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

	// Create consumer group ID
	groupID := k.config.KafkaConsumerGroup
	if groupID == "" {
		groupID = fmt.Sprintf("evero-consumer-%s", topic)
	}
	if options.QueueName != "" {
		groupID = options.QueueName
	}

	// Create consumer group configuration
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_8_0_0
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Group.Session.Timeout = 10 * time.Second
	saramaConfig.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Authentication (reuse from main config)
	if k.config.Username != "" && k.config.Password != "" {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.User = k.config.Username
		saramaConfig.Net.SASL.Password = k.config.Password

		switch k.config.KafkaSASLMechanism {
		case "SCRAM-SHA-256":
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		case "SCRAM-SHA-512":
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		case "PLAIN":
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		default:
			saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		}
	}

	if k.config.TLSEnabled {
		saramaConfig.Net.TLS.Enable = true
	}

	// Create consumer group
	brokers := k.getBrokers()
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer group: %w", err)
	}

	subCtx, cancel := context.WithCancel(ctx)

	subscription := &kafkaSubscription{
		consumerGroup: consumerGroup,
		handler:       handler,
		options:       options,
		cancel:        cancel,
		topic:         topic,
		groupID:       groupID,
	}

	k.subscribers[topic] = subscription

	// Create consumer group handler
	cgHandler := &kafkaConsumerGroupHandler{
		subscription: subscription,
		broker:       k,
	}

	// Start consuming in a goroutine
	go func() {
		defer func() {
			consumerGroup.Close()
		}()

		for {
			select {
			case <-subCtx.Done():
				return
			default:
				// Consume should be called inside an infinite loop
				if err := consumerGroup.Consume(subCtx, []string{topic}, cgHandler); err != nil {
					if errors.Is(err, sarama.ErrClosedConsumerGroup) {
						return
					}
					fmt.Printf("Error from consumer group: %v\n", err)
					time.Sleep(time.Second)
				}
			}
		}
	}()

	// Handle consumer group errors
	go func() {
		for err := range consumerGroup.Errors() {
			fmt.Printf("Consumer group error: %v\n", err)
		}
	}()

	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *kafkaConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *kafkaConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (h *kafkaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE: Do not move the code above to a goroutine
	// The `ConsumeClaim` itself is called within a goroutine
	for {
		select {
		case msg := <-claim.Messages():
			if msg == nil {
				return nil
			}
			h.handleKafkaMessage(session, msg)
		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *kafkaConsumerGroupHandler) handleKafkaMessage(session sarama.ConsumerGroupSession, kafkaMsg *sarama.ConsumerMessage) {
	message := &Message{
		ID:              fmt.Sprintf("%s-%d-%d", kafkaMsg.Topic, kafkaMsg.Partition, kafkaMsg.Offset),
		Topic:           kafkaMsg.Topic,
		Data:            kafkaMsg.Value,
		Headers:         make(map[string]string),
		Timestamp:       kafkaMsg.Timestamp,
		OriginalMessage: kafkaMsg,
	}

	// Convert headers
	for _, header := range kafkaMsg.Headers {
		message.Headers[string(header.Key)] = string(header.Value)
	}

	// Add Kafka-specific metadata
	message.Headers["kafka.partition"] = strconv.Itoa(int(kafkaMsg.Partition))
	message.Headers["kafka.offset"] = strconv.FormatInt(kafkaMsg.Offset, 10)
	if kafkaMsg.Key != nil {
		message.Headers["kafka.key"] = string(kafkaMsg.Key)
	}

	// Process message with retries
	var lastErr error
	for retry := 0; retry <= h.subscription.options.MaxRetries; retry++ {
		message.Retry = retry
		message.MaxRetries = h.subscription.options.MaxRetries

		err := h.subscription.handler(session.Context(), message)
		if err == nil {
			// Success - mark message
			session.MarkMessage(kafkaMsg, "")
			return
		}

		lastErr = err
		if retry < h.subscription.options.MaxRetries {
			time.Sleep(h.subscription.options.RetryDelay)
		}
	}

	// Failed after all retries - still mark to avoid reprocessing
	fmt.Printf("Failed to process Kafka message after %d retries: %v\n", h.subscription.options.MaxRetries, lastErr)
	session.MarkMessage(kafkaMsg, "")
}

// Unsubscribe unsubscribes from the specified topic
func (k *kafkaBroker) Unsubscribe(ctx context.Context, topic string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	subscription, exists := k.subscribers[topic]
	if !exists {
		return errSubscriptionNotFound
	}

	if subscription.cancel != nil {
		subscription.cancel()
	}

	if subscription.consumerGroup != nil {
		subscription.consumerGroup.Close()
	}

	delete(k.subscribers, topic)
	return nil
}

// CreateTopic creates a new topic in Kafka
func (k *kafkaBroker) CreateTopic(ctx context.Context, topic string, options *TopicOptions) error {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	if !k.connected {
		return errBrokerNotConnected
	}

	// Default topic configuration
	numPartitions := int32(1)
	replicationFactor := int16(1)

	if options != nil && options.Arguments != nil {
		if partitions, ok := options.Arguments["partitions"].(int); ok {
			numPartitions = int32(partitions)
		}
		if replication, ok := options.Arguments["replication-factor"].(int); ok {
			replicationFactor = int16(replication)
		}
	}

	// Create admin client from existing client
	admin, err := sarama.NewClusterAdminFromClient(k.client)
	if err != nil {
		return fmt.Errorf("failed to create Kafka admin client: %w", err)
	}
	defer admin.Close()

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	err = admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		return fmt.Errorf("failed to create Kafka topic %s: %w", topic, err)
	}

	return nil
}

// DeleteTopic deletes a topic from Kafka
func (k *kafkaBroker) DeleteTopic(ctx context.Context, topic string) error {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	if !k.connected {
		return errBrokerNotConnected
	}

	// Create admin client from existing client
	admin, err := sarama.NewClusterAdminFromClient(k.client)
	if err != nil {
		return fmt.Errorf("failed to create Kafka admin client: %w", err)
	}
	defer admin.Close()

	err = admin.DeleteTopic(topic)
	if err != nil {
		return fmt.Errorf("failed to delete Kafka topic %s: %w", topic, err)
	}

	return nil
}

// ListTopics returns a list of available topics
func (k *kafkaBroker) ListTopics(ctx context.Context) ([]string, error) {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	if !k.connected {
		return nil, errBrokerNotConnected
	}

	// Create admin client from existing client
	admin, err := sarama.NewClusterAdminFromClient(k.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka admin client: %w", err)
	}
	defer admin.Close()

	metadata, err := admin.ListTopics()
	if err != nil {
		return nil, fmt.Errorf("failed to list Kafka topics: %w", err)
	}

	topics := make([]string, 0, len(metadata))
	for topic := range metadata {
		topics = append(topics, topic)
	}

	return topics, nil
}

// Ping checks if Kafka is accessible
func (k *kafkaBroker) Ping(ctx context.Context) error {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	if !k.connected || k.client == nil {
		return errBrokerNotConnected
	}

	// Check if client is closed
	if k.client.Closed() {
		return errBrokerNotConnected
	}

	// Get brokers to test connectivity
	brokers := k.client.Brokers()
	if len(brokers) == 0 {
		return errors.New("no Kafka brokers available")
	}

	return nil
}

// PublishBatch publishes multiple messages in a batch
func (k *kafkaBroker) PublishBatch(ctx context.Context, messages []BatchMessage, options *PublishOptions) error {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	if !k.connected {
		return errBrokerNotConnected
	}

	// Convert batch messages to Sarama producer messages
	saramaMessages := make([]*sarama.ProducerMessage, 0, len(messages))

	for _, msg := range messages {
		saramaMsg := &sarama.ProducerMessage{
			Topic:     msg.Topic,
			Value:     sarama.ByteEncoder(msg.Data),
			Timestamp: time.Now(),
		}

		// Add headers
		if len(msg.Headers) > 0 {
			headers := make([]sarama.RecordHeader, 0, len(msg.Headers))
			for k, v := range msg.Headers {
				headers = append(headers, sarama.RecordHeader{
					Key:   []byte(k),
					Value: []byte(v),
				})
			}
			saramaMsg.Headers = headers
		}

		// Set key for partitioning if provided in headers
		if key, exists := msg.Headers["kafka.key"]; exists {
			saramaMsg.Key = sarama.StringEncoder(key)
		}

		saramaMessages = append(saramaMessages, saramaMsg)
	}

	// Send all messages using sync producer
	for _, msg := range saramaMessages {
		_, _, err := k.producer.SendMessage(msg)
		if err != nil {
			return fmt.Errorf("failed to publish batch message to topic %s: %w", msg.Topic, err)
		}
	}

	return nil
}

// GetStats returns broker statistics
func (k *kafkaBroker) GetStats(ctx context.Context) (*BrokerStats, error) {
	stats := &BrokerStats{
		ConnectedClients: 1, // Current connection
		Custom: map[string]interface{}{
			"type":        "Kafka",
			"brokers":     k.getBrokers(),
			"subscribers": len(k.subscribers),
		},
	}

	// Add topic-specific stats
	topics := make([]TopicStats, 0, len(k.subscribers))
	for topic := range k.subscribers {
		topics = append(topics, TopicStats{
			Name:        topic,
			Subscribers: 1,
		})
	}
	stats.Topics = topics

	return stats, nil
}

// Close closes the Kafka broker
func (k *kafkaBroker) Close() error {
	return k.Disconnect(context.Background())
}
