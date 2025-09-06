package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	messagebroker "github.com/prayaspoudel/infrastructure/message-broker"
)

func main() {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Demonstrate all three message brokers
	fmt.Println("Message Broker Demo")
	fmt.Println("===================")

	// 1. Kafka Example
	fmt.Println("\n1. Kafka Example:")
	demonstrateKafka(ctx)

	// 2. RabbitMQ Example
	fmt.Println("\n2. RabbitMQ Example:")
	demonstrateRabbitMQ(ctx)

	// 3. NATS Example
	fmt.Println("\n3. NATS Example:")
	demonstrateNATS(ctx)

	fmt.Println("\nDemo completed!")
}

func demonstrateKafka(ctx context.Context) {
	config := &messagebroker.BrokerConfig{
		KafkaBrokers:       []string{"localhost:9092"},
		KafkaConsumerGroup: "demo-group",
	}

	broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceKafka, config)
	if err != nil {
		log.Printf("Failed to create Kafka broker: %v", err)
		return
	}

	err = broker.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to Kafka: %v (Make sure Kafka is running)", err)
		return
	}
	defer broker.Close()

	fmt.Println("✓ Connected to Kafka")

	// Create topic
	topicOptions := messagebroker.KafkaTopicOptions(3, 1)
	err = broker.CreateTopic(ctx, "demo.events", topicOptions)
	if err != nil {
		log.Printf("Note: Topic creation failed (might already exist): %v", err)
	}

	// Publish messages
	publishOptions := messagebroker.KafkaPublishOptions("user-123", 0)
	err = broker.Publish(ctx, "demo.events", []byte("Hello from Kafka!"), publishOptions)
	if err != nil {
		log.Printf("Failed to publish: %v", err)
		return
	}
	fmt.Println("✓ Published message to Kafka")

	// Publish JSON message
	data := map[string]interface{}{
		"user_id":   123,
		"action":    "demo",
		"timestamp": time.Now(),
		"message":   "Hello from Kafka JSON!",
	}
	err = broker.PublishJSON(ctx, "demo.events", data, publishOptions)
	if err != nil {
		log.Printf("Failed to publish JSON: %v", err)
		return
	}
	fmt.Println("✓ Published JSON message to Kafka")

	// Subscribe and handle messages
	handler := func(ctx context.Context, msg *messagebroker.Message) error {
		fmt.Printf("📨 Kafka received: %s (partition: %s, offset: %s)\n",
			string(msg.Data),
			msg.Headers["kafka.partition"],
			msg.Headers["kafka.offset"])
		return nil
	}

	subscribeOptions := messagebroker.KafkaSubscribeOptions("demo-group", 1)
	err = broker.Subscribe(ctx, "demo.events", handler, subscribeOptions)
	if err != nil {
		log.Printf("Failed to subscribe: %v", err)
		return
	}
	fmt.Println("✓ Subscribed to Kafka topic")

	// Wait a bit for message processing
	time.Sleep(2 * time.Second)
}

func demonstrateRabbitMQ(ctx context.Context) {
	config := &messagebroker.BrokerConfig{
		RabbitMQURL:      "amqp://localhost:5672",
		RabbitMQExchange: "demo-exchange",
	}

	broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceRabbitMQ, config)
	if err != nil {
		log.Printf("Failed to create RabbitMQ broker: %v", err)
		return
	}

	err = broker.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v (Make sure RabbitMQ is running)", err)
		return
	}
	defer broker.Close()

	fmt.Println("✓ Connected to RabbitMQ")

	// Create topic/queue
	topicOptions := messagebroker.DefaultTopicOptions()
	err = broker.CreateTopic(ctx, "demo.queue", topicOptions)
	if err != nil {
		log.Printf("Note: Queue creation failed: %v", err)
	}

	// Publish messages
	publishOptions := messagebroker.JSONPublishOptions()
	err = broker.Publish(ctx, "demo.queue", []byte("Hello from RabbitMQ!"), publishOptions)
	if err != nil {
		log.Printf("Failed to publish: %v", err)
		return
	}
	fmt.Println("✓ Published message to RabbitMQ")

	// Publish JSON message
	data := map[string]interface{}{
		"user_id":   456,
		"action":    "demo",
		"timestamp": time.Now(),
		"message":   "Hello from RabbitMQ JSON!",
	}
	err = broker.PublishJSON(ctx, "demo.queue", data, publishOptions)
	if err != nil {
		log.Printf("Failed to publish JSON: %v", err)
		return
	}
	fmt.Println("✓ Published JSON message to RabbitMQ")

	// Subscribe and handle messages
	handler := func(ctx context.Context, msg *messagebroker.Message) error {
		fmt.Printf("📨 RabbitMQ received: %s (ID: %s)\n", string(msg.Data), msg.ID)
		return nil
	}

	subscribeOptions := &messagebroker.SubscribeOptions{
		QueueName:   "demo.queue",
		Durable:     true,
		AutoAck:     false,
		MaxRetries:  3,
		RetryDelay:  time.Second,
		Concurrency: 1,
	}
	err = broker.Subscribe(ctx, "demo.queue", handler, subscribeOptions)
	if err != nil {
		log.Printf("Failed to subscribe: %v", err)
		return
	}
	fmt.Println("✓ Subscribed to RabbitMQ queue")

	// Wait a bit for message processing
	time.Sleep(2 * time.Second)
}

func demonstrateNATS(ctx context.Context) {
	config := &messagebroker.BrokerConfig{
		NATSURL: "nats://localhost:4222",
	}

	broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceNATS, config)
	if err != nil {
		log.Printf("Failed to create NATS broker: %v", err)
		return
	}

	err = broker.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to NATS: %v (Make sure NATS is running)", err)
		return
	}
	defer broker.Close()

	fmt.Println("✓ Connected to NATS")

	// Publish messages
	publishOptions := messagebroker.JSONPublishOptions()
	err = broker.Publish(ctx, "demo.subject", []byte("Hello from NATS!"), publishOptions)
	if err != nil {
		log.Printf("Failed to publish: %v", err)
		return
	}
	fmt.Println("✓ Published message to NATS")

	// Publish JSON message
	data := map[string]interface{}{
		"user_id":   789,
		"action":    "demo",
		"timestamp": time.Now(),
		"message":   "Hello from NATS JSON!",
	}
	err = broker.PublishJSON(ctx, "demo.subject", data, publishOptions)
	if err != nil {
		log.Printf("Failed to publish JSON: %v", err)
		return
	}
	fmt.Println("✓ Published JSON message to NATS")

	// Subscribe and handle messages
	handler := func(ctx context.Context, msg *messagebroker.Message) error {
		fmt.Printf("📨 NATS received: %s (subject: %s)\n", string(msg.Data), msg.Topic)
		return nil
	}

	subscribeOptions := messagebroker.WorkerSubscribeOptions(1)
	subscribeOptions.QueueName = "demo-workers" // Queue group for load balancing
	err = broker.Subscribe(ctx, "demo.subject", handler, subscribeOptions)
	if err != nil {
		log.Printf("Failed to subscribe: %v", err)
		return
	}
	fmt.Println("✓ Subscribed to NATS subject")

	// Wait a bit for message processing
	time.Sleep(2 * time.Second)
}

// Example of using circuit breaker
func demonstrateCircuitBreaker() {
	// Create a handler that might fail
	unreliableHandler := func(ctx context.Context, msg *messagebroker.Message) error {
		// Simulate random failures
		if time.Now().UnixNano()%3 == 0 {
			return fmt.Errorf("simulated failure")
		}
		fmt.Printf("Successfully processed: %s\n", string(msg.Data))
		return nil
	}

	// Wrap with circuit breaker
	cb := messagebroker.NewCircuitBreakerMessageHandler(unreliableHandler, 3, 30*time.Second)

	// Use the circuit breaker handler
	_ = cb // Use cb.Handle(ctx, msg) in real usage
}

// Example of using JSON handler helper
func demonstrateJSONHandler() {
	// Define a struct for your JSON data
	type UserEvent struct {
		UserID    int       `json:"user_id"`
		Action    string    `json:"action"`
		Timestamp time.Time `json:"timestamp"`
	}

	// Create a typed handler
	typedHandler := func(ctx context.Context, event UserEvent) error {
		fmt.Printf("User %d performed %s at %v\n", event.UserID, event.Action, event.Timestamp)
		return nil
	}

	// Convert to message handler
	handler := messagebroker.JSONMessageHandler(typedHandler)

	// Use the handler in subscription
	_ = handler // Use in broker.Subscribe(ctx, topic, handler, options)
}
