# Evero Infrastructure Guide

## Overview

This guide covers the shared infrastructure components used across all Evero modules. These components are located in the `infrastructure/` directory and provide common functionality for database access, caching, logging, message brokering, routing, and validation.

## Infrastructure Components

### 1. Database (`infrastructure/database/`)

PostgreSQL database connection and query utilities.

#### Features
- Connection pooling
- Transaction management
- Query builder helpers
- Migration support
- Multi-database support (one per module)

#### Usage Example

```go
package main

import (
    "github.com/evero/infrastructure/database"
    "github.com/evero/infrastructure/config"
)

func main() {
    cfg := config.Load()
    
    // Initialize database connection
    db, err := database.NewConnection(cfg.Database)
    if err != nil {
        panic(err)
    }
    defer db.Close()
    
    // Execute query
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
    if err != nil {
        panic(err)
    }
}
```

#### Configuration

```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "database": "evero_healthcare",
    "user": "evero_user",
    "password": "secure_password",
    "sslmode": "disable",
    "max_open_conns": 25,
    "max_idle_conns": 5,
    "conn_max_lifetime": "5m"
  }
}
```

#### Connection Pooling

The database package automatically manages connection pooling:
- **max_open_conns**: Maximum number of open connections
- **max_idle_conns**: Maximum number of idle connections
- **conn_max_lifetime**: Maximum lifetime of a connection

### 2. Cache (`infrastructure/cache/`)

Redis-based caching layer for improved performance.

#### Features
- Key-value storage
- TTL (Time-To-Live) support
- Cache invalidation
- Namespace isolation per module
- Distributed caching

#### Usage Example

```go
import "github.com/evero/infrastructure/cache"

// Initialize cache
cache, err := cache.NewRedisCache(cfg.Redis)
if err != nil {
    panic(err)
}

// Set value with TTL
err = cache.Set("user:123", userData, 15*time.Minute)

// Get value
var user User
err = cache.Get("user:123", &user)

// Delete value
err = cache.Delete("user:123")

// Clear all cache with prefix
err = cache.ClearPrefix("user:")
```

#### Configuration

```json
{
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0,
    "pool_size": 10,
    "namespace": "healthcare"
  }
}
```

#### Caching Strategies

**Cache-Aside Pattern**
```go
func GetUser(id int) (*User, error) {
    // Try cache first
    var user User
    err := cache.Get(fmt.Sprintf("user:%d", id), &user)
    if err == nil {
        return &user, nil
    }
    
    // Cache miss, fetch from database
    user, err = repo.GetUser(id)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    cache.Set(fmt.Sprintf("user:%d", id), user, 15*time.Minute)
    return &user, nil
}
```

### 3. Logger (`infrastructure/logger/`)

Structured logging with multiple output formats and levels.

#### Features
- Multiple log levels (DEBUG, INFO, WARN, ERROR, FATAL)
- JSON and text output formats
- Context propagation
- Request ID tracking
- File and console output

#### Usage Example

```go
import "github.com/evero/infrastructure/logger"

// Initialize logger
log := logger.New(cfg.Logging)

// Basic logging
log.Info("Application started")
log.Debug("Debug information", "key", "value")
log.Error("Error occurred", "error", err)

// With context
ctx := logger.WithContext(context.Background(), "request_id", "abc123")
logger.InfoContext(ctx, "Processing request")

// With fields
log.With("user_id", 123).Info("User logged in")
```

#### Configuration

```json
{
  "logging": {
    "level": "info",
    "format": "json",
    "output": "stdout",
    "file_path": "/var/log/evero/healthcare.log",
    "max_size": 100,
    "max_backups": 3,
    "max_age": 28
  }
}
```

#### Log Levels

- **DEBUG**: Detailed information for debugging
- **INFO**: General informational messages
- **WARN**: Warning messages
- **ERROR**: Error messages
- **FATAL**: Critical errors that cause application termination

### 4. Message Broker (`infrastructure/message-broker/`)

RabbitMQ/Kafka integration for asynchronous communication.

#### Features
- Publish/Subscribe pattern
- Message queues
- Dead letter queues
- Message persistence
- Retry mechanisms

#### Usage Example

**Publishing Events**
```go
import "github.com/evero/infrastructure/message-broker"

// Initialize broker
broker, err := message_broker.NewRabbitMQ(cfg.RabbitMQ)
if err != nil {
    panic(err)
}

// Publish event
event := Event{
    Type: "patient.created",
    Data: patientData,
    Timestamp: time.Now(),
}
err = broker.Publish("healthcare.events", event)
```

**Consuming Events**
```go
// Subscribe to events
err = broker.Subscribe("healthcare.events", func(msg Message) error {
    var event Event
    err := json.Unmarshal(msg.Body, &event)
    if err != nil {
        return err
    }
    
    // Process event
    return processEvent(event)
})
```

#### Configuration

```json
{
  "rabbitmq": {
    "host": "localhost",
    "port": 5672,
    "user": "evero",
    "password": "password",
    "vhost": "/",
    "exchange": "evero.events",
    "queue_prefix": "healthcare"
  }
}
```

### 5. Router (`infrastructure/router/`)

HTTP router setup with middleware support using Gin framework.

#### Features
- Route registration
- Middleware chain
- Parameter binding
- Route groups
- CORS support

#### Usage Example

```go
import "github.com/evero/infrastructure/router"

// Initialize router
r := router.New(cfg.Server)

// Register middleware
r.Use(logger.Middleware())
r.Use(auth.Middleware())

// Register routes
r.GET("/health", healthHandler)
r.POST("/api/v1/patients", createPatientHandler)

// Route groups
api := r.Group("/api/v1")
{
    api.GET("/patients", listPatientsHandler)
    api.GET("/patients/:id", getPatientHandler)
    api.PUT("/patients/:id", updatePatientHandler)
    api.DELETE("/patients/:id", deletePatientHandler)
}

// Start server
r.Run(":3001")
```

#### Middleware

**Built-in Middleware**
- Logger middleware
- Recovery middleware (panic recovery)
- CORS middleware
- Authentication middleware
- Rate limiting middleware

**Custom Middleware Example**
```go
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```

### 6. Validator (`infrastructure/validator/`)

Request validation using go-playground/validator.

#### Features
- Struct validation
- Custom validators
- Error message formatting
- Automatic binding
- Field-level validation

#### Usage Example

**Model with Validation Tags**
```go
type CreatePatientRequest struct {
    FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
    LastName    string `json:"last_name" validate:"required,min=2,max=50"`
    Email       string `json:"email" validate:"required,email"`
    DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
    Phone       string `json:"phone" validate:"required,e164"`
}
```

**Validation in Handler**
```go
import "github.com/evero/infrastructure/validator"

func createPatientHandler(c *gin.Context) {
    var req CreatePatientRequest
    
    // Bind and validate
    if err := validator.BindAndValidate(c, &req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Process request...
}
```

#### Custom Validators

```go
// Register custom validator
validator.RegisterValidation("future_date", func(fl validator.FieldLevel) bool {
    date, err := time.Parse("2006-01-02", fl.Field().String())
    if err != nil {
        return false
    }
    return date.After(time.Now())
})

// Use in struct
type AppointmentRequest struct {
    Date string `validate:"required,future_date"`
}
```

### 7. Configuration (`infrastructure/config/`)

Configuration management for all modules.

#### Features
- Multiple environment support (dev, staging, production)
- JSON configuration files
- Environment variable override
- Type-safe configuration
- Hot reload support

#### Usage Example

```go
import "github.com/evero/infrastructure/config"

// Load configuration
cfg := config.Load()

// Access configuration
dbHost := cfg.Database.Host
serverPort := cfg.Server.Port
logLevel := cfg.Logging.Level

// Environment-specific config
env := os.Getenv("ENVIRONMENT")
cfg := config.LoadEnvironment(env) // dev, staging, production
```

#### Configuration Files

```
config/
├── healthcare/
│   ├── development.json
│   ├── staging.json
│   └── production.json
├── insurance/
│   ├── development.json
│   ├── staging.json
│   └── production.json
└── ...
```

#### Environment Variable Override

```bash
# Override database host
export DATABASE_HOST=db.example.com

# Override server port
export SERVER_PORT=8080

# Override log level
export LOG_LEVEL=debug
```

### 8. Setup (`infrastructure/setup/`)

Infrastructure initialization and setup utilities.

#### Features
- Database connection setup
- Migration runner
- Cache initialization
- Logger setup
- Health checks

#### Usage Example

```go
import "github.com/evero/infrastructure/setup"

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize all infrastructure
    infra, err := setup.Initialize(cfg)
    if err != nil {
        panic(err)
    }
    defer infra.Cleanup()
    
    // Run migrations
    err = setup.RunMigrations(infra.DB, "database/healthcare/migrations")
    if err != nil {
        panic(err)
    }
    
    // Start application
    startApplication(infra)
}
```

## Best Practices

### 1. Database

- **Use connection pooling**: Don't create new connections for each request
- **Close connections**: Always defer connection cleanup
- **Use transactions**: For multi-step operations
- **Prepare statements**: For repeated queries
- **Index properly**: Add indexes for frequently queried columns

### 2. Cache

- **Set appropriate TTLs**: Balance freshness vs performance
- **Use namespaces**: Separate cache keys per module
- **Handle cache misses**: Always have fallback to database
- **Invalidate on updates**: Clear cache when data changes
- **Monitor hit rate**: Track cache effectiveness

### 3. Logging

- **Use structured logging**: JSON format for easy parsing
- **Include context**: Request ID, user ID, etc.
- **Choose appropriate levels**: DEBUG for development, INFO+ for production
- **Don't log sensitive data**: Passwords, tokens, PII
- **Rotate logs**: Prevent disk space issues

### 4. Message Broker

- **Idempotent consumers**: Handle duplicate messages
- **Error handling**: Use dead letter queues
- **Message persistence**: For critical events
- **Monitoring**: Track queue depth and consumer lag
- **Versioning**: Version your event schemas

### 5. Routing

- **Group related routes**: Use route groups for organization
- **Version your API**: Use `/api/v1/`, `/api/v2/`
- **Use middleware wisely**: Apply only where needed
- **Handle errors consistently**: Use error middleware
- **Document endpoints**: Use Swagger/OpenAPI

### 6. Validation

- **Validate at boundaries**: Handler level validation
- **Custom validators**: For domain-specific rules
- **Clear error messages**: Help clients fix issues
- **Sanitize input**: Prevent injection attacks
- **Validate before processing**: Fail fast

## Performance Considerations

### Database
- Use read replicas for read-heavy workloads
- Implement query result caching
- Optimize complex queries
- Monitor slow query log

### Cache
- Use Redis cluster for high availability
- Implement cache warming for critical data
- Monitor memory usage
- Use appropriate eviction policies

### Logging
- Use asynchronous logging for high-throughput
- Compress old log files
- Send logs to centralized system (ELK, Splunk)
- Sample verbose logs in production

### Message Broker
- Use message batching for high throughput
- Implement backpressure mechanisms
- Monitor consumer lag
- Scale consumers horizontally

## Monitoring

### Health Checks

```go
func healthCheck(db *sql.DB, cache *redis.Client) error {
    // Check database
    if err := db.Ping(); err != nil {
        return fmt.Errorf("database unhealthy: %w", err)
    }
    
    // Check cache
    if err := cache.Ping().Err(); err != nil {
        return fmt.Errorf("cache unhealthy: %w", err)
    }
    
    return nil
}
```

### Metrics

Expose metrics for monitoring:
- Request rate and latency
- Database connection pool stats
- Cache hit/miss ratio
- Message queue depth
- Error rates

## Troubleshooting

### Database Connection Issues
```bash
# Check connection
psql -U evero_user -d evero_healthcare

# View active connections
SELECT * FROM pg_stat_activity;

# Check connection pool stats
SELECT * FROM pg_stat_database;
```

### Redis Connection Issues
```bash
# Test connection
redis-cli ping

# Monitor commands
redis-cli monitor

# Check memory usage
redis-cli info memory
```

### RabbitMQ Issues
```bash
# List queues
rabbitmqctl list_queues

# Check connection count
rabbitmqctl list_connections

# View queue details
rabbitmqctl list_queues name messages consumers
```

## Security

### Database Security
- Use least privilege principle
- Encrypt connections (SSL/TLS)
- Rotate credentials regularly
- Use prepared statements (prevent SQL injection)

### Cache Security
- Set password for Redis
- Use SSL for connections
- Implement access control
- Don't cache sensitive data in plaintext

### API Security
- Use HTTPS only
- Implement rate limiting
- Validate all input
- Use CORS properly
- Implement authentication/authorization

## Example Usage

See `infrastructure/example_usage.md` for comprehensive examples of using all infrastructure components together.

## Conclusion

The Evero infrastructure provides a solid foundation for building scalable, maintainable applications. By leveraging these shared components, modules can focus on business logic while relying on battle-tested infrastructure.

For module-specific implementations, refer to the respective module documentation in `docs/[module]/`.
