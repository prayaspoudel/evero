// Package cache provides caching functionality for storing and retrieving
// data from various cache backends including Redis and in-memory stores.
package cache

import (
	"context"
	"time"
)

// CacheManager interface defines the contract for cache management
type CacheManager interface {
	// Connect establishes connection to the cache backend
	Connect(ctx context.Context) error

	// Disconnect closes the connection to the cache backend
	Disconnect(ctx context.Context) error

	// Set stores a value with the given key and expiration time
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get retrieves a value by key
	Get(ctx context.Context, key string) (interface{}, error)

	// GetString retrieves a string value by key
	GetString(ctx context.Context, key string) (string, error)

	// GetInt retrieves an integer value by key
	GetInt(ctx context.Context, key string) (int, error)

	// GetBool retrieves a boolean value by key
	GetBool(ctx context.Context, key string) (bool, error)

	// GetFloat64 retrieves a float64 value by key
	GetFloat64(ctx context.Context, key string) (float64, error)

	// Delete removes a value by key
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in the cache
	Exists(ctx context.Context, key string) (bool, error)

	// Keys returns all keys matching the given pattern
	Keys(ctx context.Context, pattern string) ([]string, error)

	// Expire sets an expiration time for a key
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// TTL returns the time to live for a key
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Clear removes all keys from the cache
	Clear(ctx context.Context) error

	// Ping checks if the cache backend is accessible
	Ping(ctx context.Context) error

	// SetMultiple stores multiple key-value pairs
	SetMultiple(ctx context.Context, pairs map[string]interface{}, expiration time.Duration) error

	// GetMultiple retrieves multiple values by keys
	GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error)

	// DeleteMultiple removes multiple keys
	DeleteMultiple(ctx context.Context, keys []string) error

	// Increment increments a numeric value
	Increment(ctx context.Context, key string, value int64) (int64, error)

	// Decrement decrements a numeric value
	Decrement(ctx context.Context, key string, value int64) (int64, error)

	// Close closes the cache manager and releases resources
	Close() error
}

// CacheConfig holds configuration for cache backends
type CacheConfig struct {
	// Redis configuration
	RedisAddr     string `json:"redis_addr"`
	RedisPassword string `json:"redis_password"`
	RedisDB       int    `json:"redis_db"`

	// Connection pool settings
	MaxRetries   int           `json:"max_retries"`
	PoolSize     int           `json:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	PoolTimeout  time.Duration `json:"pool_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`

	// In-memory cache settings
	DefaultExpiration time.Duration `json:"default_expiration"`
	CleanupInterval   time.Duration `json:"cleanup_interval"`
	MaxSize           int           `json:"max_size"`
}
