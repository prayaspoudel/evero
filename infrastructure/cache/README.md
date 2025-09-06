# Cache Management

This package provides a unified interface for cache management with support for multiple backends including Redis and in-memory caching.

## Features

- **Multiple Backends**: Support for Redis and in-memory caching
- **Unified Interface**: Consistent API across all cache backends
- **Context Support**: All operations support context for cancellation and timeouts
- **Expiration Management**: Support for key expiration and TTL
- **Batch Operations**: Support for batch set/get/delete operations
- **Atomic Operations**: Support for increment/decrement operations
- **Pattern Matching**: Support for key pattern matching
- **Connection Management**: Proper connection handling and cleanup

## Supported Backends

### Redis
- Full Redis functionality including clustering support
- Connection pooling and retry mechanisms
- Persistent storage across application restarts

### In-Memory
- Fast local caching with automatic cleanup
- Memory-efficient with configurable size limits
- No external dependencies

## Usage

### Basic Usage

```go
import "your-project/infrastructure/cache"

// Create Redis cache manager
config := &cache.CacheConfig{
    RedisAddr:     "localhost:6379",
    RedisPassword: "",
    RedisDB:       0,
    PoolSize:      10,
}

cacheManager, err := cache.NewCacheManagerFactory(cache.InstanceRedis, config)
if err != nil {
    log.Fatal(err)
}

// Connect to cache
ctx := context.Background()
err = cacheManager.Connect(ctx)
if err != nil {
    log.Fatal(err)
}
defer cacheManager.Close()

// Set a value
err = cacheManager.Set(ctx, "key", "value", 5*time.Minute)
if err != nil {
    log.Fatal(err)
}

// Get a value
value, err := cacheManager.GetString(ctx, "key")
if err != nil {
    log.Fatal(err)
}
fmt.Println(value)
```

### In-Memory Cache

```go
// Create in-memory cache manager
config := &cache.CacheConfig{
    DefaultExpiration: 5 * time.Minute,
    CleanupInterval:   10 * time.Minute,
    MaxSize:           1000,
}

cacheManager, err := cache.NewCacheManagerFactory(cache.InstanceInMemory, config)
if err != nil {
    log.Fatal(err)
}

// Connect (starts cleanup goroutine)
ctx := context.Background()
err = cacheManager.Connect(ctx)
if err != nil {
    log.Fatal(err)
}
defer cacheManager.Close()

// Use the cache
err = cacheManager.Set(ctx, "key", "value", time.Hour)
// ... rest of operations
```

### Batch Operations

```go
// Set multiple values
pairs := map[string]interface{}{
    "key1": "value1",
    "key2": 42,
    "key3": true,
}
err = cacheManager.SetMultiple(ctx, pairs, 10*time.Minute)

// Get multiple values
keys := []string{"key1", "key2", "key3"}
values, err := cacheManager.GetMultiple(ctx, keys)

// Delete multiple keys
err = cacheManager.DeleteMultiple(ctx, keys)
```

### Atomic Operations

```go
// Increment a counter
newValue, err := cacheManager.Increment(ctx, "counter", 1)

// Decrement a counter
newValue, err := cacheManager.Decrement(ctx, "counter", 1)
```

## Configuration

### Redis Configuration

```go
type CacheConfig struct {
    // Redis settings
    RedisAddr     string        `json:"redis_addr"`     // Redis server address
    RedisPassword string        `json:"redis_password"` // Redis password
    RedisDB       int           `json:"redis_db"`       // Redis database number
    
    // Connection pool settings
    MaxRetries     int           `json:"max_retries"`     // Maximum retry attempts
    PoolSize       int           `json:"pool_size"`       // Connection pool size
    MinIdleConns   int           `json:"min_idle_conns"`  // Minimum idle connections
    DialTimeout    time.Duration `json:"dial_timeout"`    // Connection timeout
    ReadTimeout    time.Duration `json:"read_timeout"`    // Read timeout
    WriteTimeout   time.Duration `json:"write_timeout"`   // Write timeout
    PoolTimeout    time.Duration `json:"pool_timeout"`    // Pool timeout
    IdleTimeout    time.Duration `json:"idle_timeout"`    // Idle connection timeout
}
```

### In-Memory Configuration

```go
type CacheConfig struct {
    // In-memory settings
    DefaultExpiration time.Duration `json:"default_expiration"` // Default expiration time
    CleanupInterval   time.Duration `json:"cleanup_interval"`   // Cleanup interval
    MaxSize           int           `json:"max_size"`           // Maximum number of items
}
```

## Error Handling

The package defines several error types:

- `errInvalidCacheInstance`: Invalid cache backend type
- `errCacheNotConnected`: Cache manager not connected
- `errKeyNotFound`: Key not found in cache
- `errInvalidKeyType`: Invalid type for key operation

## Dependencies

### Redis Backend
```
go get github.com/redis/go-redis/v9
```

### In-Memory Backend
No external dependencies required.

## Best Practices

1. **Always use context**: Pass appropriate context for timeouts and cancellation
2. **Handle errors**: Check for `errKeyNotFound` when getting values
3. **Set appropriate TTL**: Use reasonable expiration times to prevent memory leaks
4. **Close connections**: Always call `Close()` when done with cache manager
5. **Use batch operations**: For multiple operations, use batch methods for better performance
6. **Monitor memory usage**: For in-memory cache, monitor the `MaxSize` setting
