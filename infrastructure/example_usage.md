# Infrastructure Integration Example

This example shows how to integrate the cache and message-broker infrastructure components into your application.

## Required Dependencies

Add these dependencies to your `go.mod` file:

```bash
# For Redis cache support
go get github.com/redis/go-redis/v9

# For RabbitMQ message broker support  
go get github.com/streadway/amqp

# For NATS message broker support
go get github.com/nats-io/nats.go
```

## Example Usage

### 1. Cache Integration

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/prayaspoudel/infrastructure/cache"
)

func setupCache() cache.CacheManager {
    // For Redis
    config := &cache.CacheConfig{
        RedisAddr:     "localhost:6379",
        RedisPassword: "",
        RedisDB:       0,
        PoolSize:      10,
        MaxRetries:    3,
        DialTimeout:   5 * time.Second,
        ReadTimeout:   3 * time.Second,
        WriteTimeout:  3 * time.Second,
    }

    cacheManager, err := cache.NewCacheManagerFactory(cache.InstanceRedis, config)
    if err != nil {
        log.Fatal("Failed to create cache manager:", err)
    }

    // For In-Memory (alternative)
    // config := &cache.CacheConfig{
    //     DefaultExpiration: 5 * time.Minute,
    //     CleanupInterval:   10 * time.Minute,
    //     MaxSize:           1000,
    // }
    // cacheManager, err := cache.NewCacheManagerFactory(cache.InstanceInMemory, config)

    ctx := context.Background()
    if err := cacheManager.Connect(ctx); err != nil {
        log.Fatal("Failed to connect to cache:", err)
    }

    return cacheManager
}

func cacheExample() {
    cacheManager := setupCache()
    defer cacheManager.Close()

    ctx := context.Background()

    // Set a value
    err := cacheManager.Set(ctx, "user:123", "John Doe", time.Hour)
    if err != nil {
        log.Printf("Failed to set cache: %v", err)
        return
    }

    // Get a value
    value, err := cacheManager.GetString(ctx, "user:123")
    if err != nil {
        log.Printf("Failed to get cache: %v", err)
        return
    }
    log.Printf("Cached value: %s", value)

    // Batch operations
    pairs := map[string]interface{}{
        "user:124": "Jane Doe",
        "user:125": "Bob Smith",
    }
    err = cacheManager.SetMultiple(ctx, pairs, time.Hour)
    if err != nil {
        log.Printf("Failed to set multiple: %v", err)
    }

    // Increment counter
    count, err := cacheManager.Increment(ctx, "page_views", 1)
    if err != nil {
        log.Printf("Failed to increment: %v", err)
    } else {
        log.Printf("Page views: %d", count)
    }
}
```

### 2. Message Broker Integration

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/prayaspoudel/infrastructure/message-broker"
)

func setupMessageBroker() messagebroker.MessageBroker {
    // For RabbitMQ
    config := &messagebroker.BrokerConfig{
        RabbitMQURL:      "amqp://localhost:5672",
        RabbitMQExchange: "evero-exchange",
        MaxReconnects:    3,
        ReconnectWait:    5 * time.Second,
        Timeout:          10 * time.Second,
    }

    broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceRabbitMQ, config)
    if err != nil {
        log.Fatal("Failed to create message broker:", err)
    }

    // For NATS (alternative)
    // config := &messagebroker.BrokerConfig{
    //     NATSURL:       "nats://localhost:4222",
    //     MaxReconnects: 3,
    //     ReconnectWait: 5 * time.Second,
    //     PingInterval:  2 * time.Minute,
    // }
    // broker, err := messagebroker.NewMessageBrokerFactory(messagebroker.InstanceNATS, config)

    ctx := context.Background()
    if err := broker.Connect(ctx); err != nil {
        log.Fatal("Failed to connect to message broker:", err)
    }

    return broker
}

func messageBrokerExample() {
    broker := setupMessageBroker()
    defer broker.Close()

    ctx := context.Background()

    // Subscribe to user events
    userEventHandler := func(ctx context.Context, msg *messagebroker.Message) error {
        log.Printf("Received user event: %s", string(msg.Data))
        
        var event map[string]interface{}
        if err := json.Unmarshal(msg.Data, &event); err != nil {
            return err
        }
        
        // Process the event
        log.Printf("Processing event: %+v", event)
        return nil
    }

    subscribeOptions := &messagebroker.SubscribeOptions{
        QueueName:     "user-events-queue",
        Durable:       true,
        AutoAck:       false,
        MaxRetries:    3,
        RetryDelay:    5 * time.Second,
        Concurrency:   2,
        PrefetchCount: 5,
    }

    err := broker.Subscribe(ctx, "user.events", userEventHandler, subscribeOptions)
    if err != nil {
        log.Printf("Failed to subscribe: %v", err)
        return
    }

    // Publish a user event
    userEvent := map[string]interface{}{
        "user_id":   123,
        "action":    "login",
        "timestamp": time.Now(),
        "ip":        "192.168.1.1",
    }

    publishOptions := &messagebroker.PublishOptions{
        Persistent: true,
        Headers: map[string]string{
            "event_type": "user_action",
            "version":    "1.0",
        },
    }

    err = broker.PublishJSON(ctx, "user.events", userEvent, publishOptions)
    if err != nil {
        log.Printf("Failed to publish event: %v", err)
        return
    }

    // Batch publish multiple events
    events := []messagebroker.BatchMessage{
        {
            Topic: "user.events",
            Data:  []byte(`{"user_id": 124, "action": "logout"}`),
            Headers: map[string]string{"event_type": "user_action"},
        },
        {
            Topic: "system.events",
            Data:  []byte(`{"component": "auth", "status": "healthy"}`),
            Headers: map[string]string{"event_type": "health_check"},
        },
    }

    err = broker.PublishBatch(ctx, events, publishOptions)
    if err != nil {
        log.Printf("Failed to publish batch: %v", err)
    }
}
```

### 3. Combined Application Example

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/prayaspoudel/infrastructure/cache"
    "github.com/prayaspoudel/infrastructure/message-broker"
)

type UserService struct {
    cache  cache.CacheManager
    broker messagebroker.MessageBroker
}

func NewUserService(cacheManager cache.CacheManager, broker messagebroker.MessageBroker) *UserService {
    return &UserService{
        cache:  cacheManager,
        broker: broker,
    }
}

func (s *UserService) CreateUser(ctx context.Context, user User) error {
    // Save user to database (not shown)
    
    // Cache the user
    userKey := fmt.Sprintf("user:%d", user.ID)
    userData, _ := json.Marshal(user)
    err := s.cache.Set(ctx, userKey, userData, time.Hour)
    if err != nil {
        log.Printf("Failed to cache user: %v", err)
    }

    // Publish user creation event
    event := map[string]interface{}{
        "event_type": "user_created",
        "user_id":    user.ID,
        "email":      user.Email,
        "timestamp":  time.Now(),
    }

    publishOptions := &messagebroker.PublishOptions{
        Persistent: true,
        Headers: map[string]string{
            "event_type": "user_lifecycle",
            "action":     "create",
        },
    }

    err = s.broker.PublishJSON(ctx, "user.lifecycle", event, publishOptions)
    if err != nil {
        log.Printf("Failed to publish user creation event: %v", err)
    }

    return nil
}

func (s *UserService) GetUser(ctx context.Context, userID int) (*User, error) {
    userKey := fmt.Sprintf("user:%d", userID)
    
    // Try to get from cache first
    cached, err := s.cache.GetString(ctx, userKey)
    if err == nil {
        var user User
        if json.Unmarshal([]byte(cached), &user) == nil {
            // Increment cache hit counter
            s.cache.Increment(ctx, "cache_hits", 1)
            return &user, nil
        }
    }

    // Cache miss - get from database (not shown)
    // user := getUserFromDB(userID)
    
    // Cache the result
    userData, _ := json.Marshal(user)
    s.cache.Set(ctx, userKey, userData, time.Hour)
    
    // Increment cache miss counter
    s.cache.Increment(ctx, "cache_misses", 1)
    
    return user, nil
}

type User struct {
    ID    int    `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

func main() {
    // Setup infrastructure
    cacheManager := setupCache()
    broker := setupMessageBroker()
    
    defer cacheManager.Close()
    defer broker.Close()

    // Setup event handlers
    setupEventHandlers(broker)

    // Create user service
    userService := NewUserService(cacheManager, broker)

    // Example usage
    ctx := context.Background()
    
    user := User{
        ID:    123,
        Email: "john@example.com",
        Name:  "John Doe",
    }

    // Create user (will cache and publish event)
    err := userService.CreateUser(ctx, user)
    if err != nil {
        log.Printf("Failed to create user: %v", err)
    }

    // Get user (will try cache first)
    retrievedUser, err := userService.GetUser(ctx, 123)
    if err != nil {
        log.Printf("Failed to get user: %v", err)
    } else {
        log.Printf("Retrieved user: %+v", retrievedUser)
    }

    // Keep the application running to process messages
    select {}
}

func setupEventHandlers(broker messagebroker.MessageBroker) {
    ctx := context.Background()

    // User lifecycle event handler
    userLifecycleHandler := func(ctx context.Context, msg *messagebroker.Message) error {
        log.Printf("Processing user lifecycle event: %s", string(msg.Data))
        
        var event map[string]interface{}
        if err := json.Unmarshal(msg.Data, &event); err != nil {
            return err
        }

        // Process the event (e.g., send welcome email, update analytics, etc.)
        eventType := event["event_type"].(string)
        userID := event["user_id"].(float64)
        
        log.Printf("Event: %s for user: %.0f", eventType, userID)
        
        return nil
    }

    subscribeOptions := messagebroker.DefaultSubscribeOptions()
    subscribeOptions.QueueName = "user-lifecycle-queue"
    subscribeOptions.Concurrency = 3

    err := broker.Subscribe(ctx, "user.lifecycle", userLifecycleHandler, subscribeOptions)
    if err != nil {
        log.Printf("Failed to subscribe to user lifecycle events: %v", err)
    }
}
```

### 4. Configuration Management

You can create configuration files for different environments:

**config/healthcare/development.json:**
```json
{
  "cache": {
    "type": "inmemory",
    "default_expiration": "5m",
    "cleanup_interval": "10m",
    "max_size": 1000
  },
  "message_broker": {
    "type": "nats",
    "nats_url": "nats://localhost:4222",
    "max_reconnects": 3,
    "reconnect_wait": "5s"
  }
}
```

**config/healthcare/production.json:**
```json
{
  "cache": {
    "type": "redis",
    "redis_addr": "prod-redis:6379",
    "redis_password": "${REDIS_PASSWORD}",
    "pool_size": 20,
    "max_retries": 5
  },
  "message_broker": {
    "type": "rabbitmq",
    "rabbitmq_url": "amqp://prod-rabbitmq:5672",
    "rabbitmq_exchange": "evero-prod",
    "max_reconnects": 5,
    "reconnect_wait": "10s"
  }
}
```

Then load the configuration using the existing config infrastructure:

```go
func loadInfrastructureConfig(environment, module string) (*InfrastructureConfig, error) {
    configManager, err := config.NewConfigManagerFactory(config.InstanceViper)
    if err != nil {
        return nil, err
    }

    err = configManager.Load(environment, module)
    if err != nil {
        return nil, err
    }

    // Create cache config
    var cacheConfig *cache.CacheConfig
    cacheType := configManager.GetString("cache.type")
    
    if cacheType == "redis" {
        cacheConfig = &cache.CacheConfig{
            RedisAddr:     configManager.GetString("cache.redis_addr"),
            RedisPassword: configManager.GetString("cache.redis_password"),
            PoolSize:      configManager.GetInt("cache.pool_size"),
            MaxRetries:    configManager.GetInt("cache.max_retries"),
        }
    } else {
        cacheConfig = &cache.CacheConfig{
            DefaultExpiration: configManager.GetDuration("cache.default_expiration"),
            CleanupInterval:   configManager.GetDuration("cache.cleanup_interval"),
            MaxSize:           configManager.GetInt("cache.max_size"),
        }
    }

    // Create broker config
    var brokerConfig *messagebroker.BrokerConfig
    brokerType := configManager.GetString("message_broker.type")
    
    if brokerType == "rabbitmq" {
        brokerConfig = &messagebroker.BrokerConfig{
            RabbitMQURL:      configManager.GetString("message_broker.rabbitmq_url"),
            RabbitMQExchange: configManager.GetString("message_broker.rabbitmq_exchange"),
            MaxReconnects:    configManager.GetInt("message_broker.max_reconnects"),
            ReconnectWait:    configManager.GetDuration("message_broker.reconnect_wait"),
        }
    } else {
        brokerConfig = &messagebroker.BrokerConfig{
            NATSURL:       configManager.GetString("message_broker.nats_url"),
            MaxReconnects: configManager.GetInt("message_broker.max_reconnects"),
            ReconnectWait: configManager.GetDuration("message_broker.reconnect_wait"),
        }
    }

    return &InfrastructureConfig{
        Cache:  cacheConfig,
        Broker: brokerConfig,
    }, nil
}

type InfrastructureConfig struct {
    Cache  *cache.CacheConfig
    Broker *messagebroker.BrokerConfig
}
```

This implementation provides:

1. **Cache Management**: Both Redis and in-memory implementations
2. **Message Broker**: Both RabbitMQ and NATS implementations  
3. **Unified Interface**: Consistent APIs across backends
4. **Configuration Integration**: Works with existing config system
5. **Production Ready**: Error handling, connection management, retries
6. **Testing Support**: Example tests and mocks
7. **Performance**: Optimized for high throughput scenarios
