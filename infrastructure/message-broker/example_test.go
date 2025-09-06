package messagebroker_test

import (
	"context"
	"testing"
	"time"

	messagebroker "github.com/prayaspoudel/infrastructure/message-broker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRabbitMQBroker(t *testing.T) {
	// Skip if RabbitMQ is not available
	t.Skip("RabbitMQ integration test - requires running RabbitMQ server")

	config := &messagebroker.BrokerConfig{
		RabbitMQURL:      "amqp://localhost:5672",
		RabbitMQExchange: "test-exchange",
	}

	broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceRabbitMQ, config)
	require.NoError(t, err)

	ctx := context.Background()
	err = broker.Connect(ctx)
	require.NoError(t, err)
	defer broker.Close()

	// Test ping
	err = broker.Ping(ctx)
	assert.NoError(t, err)

	// Test publish
	err = broker.Publish(ctx, "test.topic", []byte("Hello RabbitMQ"), nil)
	assert.NoError(t, err)

	// Test publish JSON
	data := map[string]interface{}{
		"message":   "Hello from JSON",
		"timestamp": time.Now(),
	}
	err = broker.PublishJSON(ctx, "test.json", data, messagebroker.JSONPublishOptions())
	assert.NoError(t, err)

	// Test subscribe (would need a running consumer in real test)
	handler := func(ctx context.Context, msg *messagebroker.Message) error {
		t.Logf("Received: %s", string(msg.Data))
		return nil
	}

	err = broker.Subscribe(ctx, "test.topic", handler, messagebroker.DefaultSubscribeOptions())
	assert.NoError(t, err)

	// Wait a bit for message processing
	time.Sleep(time.Second)
}

func TestNATSBroker(t *testing.T) {
	// Skip if NATS is not available
	t.Skip("NATS integration test - requires running NATS server")

	config := &messagebroker.BrokerConfig{
		NATSURL: "nats://localhost:4222",
	}

	broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceNATS, config)
	require.NoError(t, err)

	ctx := context.Background()
	err = broker.Connect(ctx)
	require.NoError(t, err)
	defer broker.Close()

	// Test ping
	err = broker.Ping(ctx)
	assert.NoError(t, err)

	// Test publish
	err = broker.Publish(ctx, "test.subject", []byte("Hello NATS"), nil)
	assert.NoError(t, err)

	// Test publish JSON
	data := map[string]interface{}{
		"message":   "Hello from NATS JSON",
		"timestamp": time.Now(),
	}
	err = broker.PublishJSON(ctx, "test.json", data, messagebroker.JSONPublishOptions())
	assert.NoError(t, err)

	// Test subscribe
	handler := func(ctx context.Context, msg *messagebroker.Message) error {
		t.Logf("Received: %s", string(msg.Data))
		return nil
	}

	options := messagebroker.WorkerSubscribeOptions(2)
	err = broker.Subscribe(ctx, "test.subject", handler, options)
	assert.NoError(t, err)

	// Wait a bit for message processing
	time.Sleep(time.Second)
}

func TestKafkaBroker(t *testing.T) {
	// Skip if Kafka is not available
	t.Skip("Kafka integration test - requires running Kafka server")

	config := &messagebroker.BrokerConfig{
		KafkaBrokers:       []string{"localhost:9092"},
		KafkaConsumerGroup: "test-group",
	}

	broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceKafka, config)
	require.NoError(t, err)

	ctx := context.Background()
	err = broker.Connect(ctx)
	require.NoError(t, err)
	defer broker.Close()

	// Test ping
	err = broker.Ping(ctx)
	assert.NoError(t, err)

	// Create topic first
	topicOptions := messagebroker.KafkaTopicOptions(3, 1) // 3 partitions, replication factor 1
	err = broker.CreateTopic(ctx, "test-topic", topicOptions)
	assert.NoError(t, err)

	// Test publish with key
	publishOptions := messagebroker.KafkaPublishOptions("user-123", 0)
	err = broker.Publish(ctx, "test-topic", []byte("Hello Kafka"), publishOptions)
	assert.NoError(t, err)

	// Test publish JSON
	data := map[string]interface{}{
		"user_id":   123,
		"action":    "login",
		"timestamp": time.Now(),
	}
	err = broker.PublishJSON(ctx, "test-topic", data, messagebroker.KafkaPublishOptions("user-123", -1))
	assert.NoError(t, err)

	// Test batch publish
	messages := []messagebroker.BatchMessage{
		{
			Topic: "test-topic",
			Data:  []byte("batch message 1"),
			Headers: map[string]string{
				"kafka.key": "batch-1",
			},
		},
		{
			Topic: "test-topic",
			Data:  []byte("batch message 2"),
			Headers: map[string]string{
				"kafka.key": "batch-2",
			},
		},
	}
	err = broker.PublishBatch(ctx, messages, nil)
	assert.NoError(t, err)

	// Test subscribe
	handler := func(ctx context.Context, msg *messagebroker.Message) error {
		t.Logf("Received from partition %s, offset %s: %s",
			msg.Headers["kafka.partition"],
			msg.Headers["kafka.offset"],
			string(msg.Data))
		return nil
	}

	subscribeOptions := messagebroker.KafkaSubscribeOptions("test-group", 2)
	err = broker.Subscribe(ctx, "test-topic", handler, subscribeOptions)
	assert.NoError(t, err)

	// Wait a bit for message processing
	time.Sleep(time.Second * 5)

	// Test list topics
	topics, err := broker.ListTopics(ctx)
	assert.NoError(t, err)
	t.Logf("Available topics: %v", topics)

	// Test get stats
	stats, err := broker.GetStats(ctx)
	assert.NoError(t, err)
	t.Logf("Broker stats: %+v", stats)
}

func TestBrokerFactory(t *testing.T) {
	tests := []struct {
		name        string
		instance    int
		config      *messagebroker.BrokerConfig
		expectError bool
	}{
		{
			name:     "RabbitMQ with valid config",
			instance: messagebroker.InstanceRabbitMQ,
			config: &messagebroker.BrokerConfig{
				RabbitMQURL: "amqp://localhost:5672",
			},
			expectError: false,
		},
		{
			name:     "NATS with valid config",
			instance: messagebroker.InstanceNATS,
			config: &messagebroker.BrokerConfig{
				NATSURL: "nats://localhost:4222",
			},
			expectError: false,
		},
		{
			name:     "Kafka with valid config",
			instance: messagebroker.InstanceKafka,
			config: &messagebroker.BrokerConfig{
				KafkaBrokers: []string{"localhost:9092"},
			},
			expectError: false,
		},
		{
			name:        "Invalid instance",
			instance:    999,
			config:      &messagebroker.BrokerConfig{},
			expectError: true,
		},
		{
			name:        "Nil config",
			instance:    messagebroker.InstanceRabbitMQ,
			config:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			broker, err := messagebroker.NewMessageBrokerFactory(tt.instance, tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, broker)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, broker)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test default options
	publishOpts := messagebroker.DefaultPublishOptions()
	assert.NotNil(t, publishOpts)
	assert.NotNil(t, publishOpts.Headers)

	subscribeOpts := messagebroker.DefaultSubscribeOptions()
	assert.NotNil(t, subscribeOpts)
	assert.Equal(t, 3, subscribeOpts.MaxRetries)

	topicOpts := messagebroker.DefaultTopicOptions()
	assert.NotNil(t, topicOpts)
	assert.True(t, topicOpts.Durable)

	// Test JSON options
	jsonOpts := messagebroker.JSONPublishOptions()
	assert.Equal(t, "application/json", jsonOpts.ContentType)
	assert.True(t, jsonOpts.Persistent)

	// Test worker options
	workerOpts := messagebroker.WorkerSubscribeOptions(5)
	assert.Equal(t, 5, workerOpts.Concurrency)
	assert.Equal(t, 10, workerOpts.PrefetchCount)

	// Test Kafka-specific options
	kafkaPublishOpts := messagebroker.KafkaPublishOptions("test-key", 2)
	assert.Equal(t, "test-key", kafkaPublishOpts.Headers["kafka.key"])
	assert.Equal(t, "2", kafkaPublishOpts.Headers["kafka.partition"])

	kafkaSubscribeOpts := messagebroker.KafkaSubscribeOptions("my-group", 3)
	assert.Equal(t, "my-group", kafkaSubscribeOpts.QueueName)
	assert.Equal(t, 3, kafkaSubscribeOpts.Concurrency)

	kafkaTopicOpts := messagebroker.KafkaTopicOptions(5, 2)
	assert.Equal(t, 5, kafkaTopicOpts.Arguments["partitions"])
	assert.Equal(t, 2, kafkaTopicOpts.Arguments["replication-factor"])
}

func TestMessageValidation(t *testing.T) {
	// Valid message
	msg := &messagebroker.Message{
		ID:        "test-id",
		Topic:     "test.topic",
		Data:      []byte("test data"),
		Headers:   map[string]string{"key": "value"},
		Timestamp: time.Now(),
	}
	err := messagebroker.ValidateMessage(msg)
	assert.NoError(t, err)

	// Nil message
	err = messagebroker.ValidateMessage(nil)
	assert.Error(t, err)

	// Empty topic
	msg.Topic = ""
	err = messagebroker.ValidateMessage(msg)
	assert.Error(t, err)

	// Nil data
	msg.Topic = "test.topic"
	msg.Data = nil
	err = messagebroker.ValidateMessage(msg)
	assert.Error(t, err)
}

func TestBatchMessageCreation(t *testing.T) {
	// Test with byte slice
	batchMsg, err := messagebroker.CreateBatchMessage("test.topic", []byte("test data"), nil)
	assert.NoError(t, err)
	assert.Equal(t, "test.topic", batchMsg.Topic)
	assert.Equal(t, []byte("test data"), batchMsg.Data)

	// Test with string
	batchMsg, err = messagebroker.CreateBatchMessage("test.topic", "test string", nil)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test string"), batchMsg.Data)

	// Test with JSON object
	data := map[string]string{"key": "value"}
	batchMsg, err = messagebroker.CreateBatchMessage("test.topic", data, map[string]string{"type": "json"})
	assert.NoError(t, err)
	assert.Contains(t, string(batchMsg.Data), "key")
	assert.Equal(t, "json", batchMsg.Headers["type"])
}
