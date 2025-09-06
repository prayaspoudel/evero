package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCacheManager struct {
	client *redis.Client
	config *CacheConfig
}

// NewRedisCacheManager creates a new Redis-based cache manager
func NewRedisCacheManager(config *CacheConfig) (CacheManager, error) {
	if config == nil {
		return nil, errors.New("cache config is required")
	}

	return &redisCacheManager{
		config: config,
	}, nil
}

// Connect establishes connection to Redis
func (r *redisCacheManager) Connect(ctx context.Context) error {
	r.client = redis.NewClient(&redis.Options{
		Addr:            r.config.RedisAddr,
		Password:        r.config.RedisPassword,
		DB:              r.config.RedisDB,
		MaxRetries:      r.config.MaxRetries,
		PoolSize:        r.config.PoolSize,
		MinIdleConns:    r.config.MinIdleConns,
		DialTimeout:     r.config.DialTimeout,
		ReadTimeout:     r.config.ReadTimeout,
		WriteTimeout:    r.config.WriteTimeout,
		PoolTimeout:     r.config.PoolTimeout,
		ConnMaxIdleTime: r.config.IdleTimeout,
	})

	return r.client.Ping(ctx).Err()
}

// Disconnect closes the Redis connection
func (r *redisCacheManager) Disconnect(ctx context.Context) error {
	if r.client == nil {
		return nil
	}
	return r.client.Close()
}

// Set stores a value with the given key and expiration time
func (r *redisCacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if r.client == nil {
		return errCacheNotConnected
	}

	var data []byte
	var err error

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a value by key
func (r *redisCacheManager) Get(ctx context.Context, key string) (interface{}, error) {
	if r.client == nil {
		return nil, errCacheNotConnected
	}

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errKeyNotFound
		}
		return nil, err
	}

	return val, nil
}

// GetString retrieves a string value by key
func (r *redisCacheManager) GetString(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", errCacheNotConnected
	}

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errKeyNotFound
		}
		return "", err
	}

	return val, nil
}

// GetInt retrieves an integer value by key
func (r *redisCacheManager) GetInt(ctx context.Context, key string) (int, error) {
	val, err := r.GetString(ctx, key)
	if err != nil {
		return 0, err
	}

	result, err := strconv.Atoi(val)
	if err != nil {
		return 0, errInvalidKeyType
	}

	return result, nil
}

// GetBool retrieves a boolean value by key
func (r *redisCacheManager) GetBool(ctx context.Context, key string) (bool, error) {
	val, err := r.GetString(ctx, key)
	if err != nil {
		return false, err
	}

	result, err := strconv.ParseBool(val)
	if err != nil {
		return false, errInvalidKeyType
	}

	return result, nil
}

// GetFloat64 retrieves a float64 value by key
func (r *redisCacheManager) GetFloat64(ctx context.Context, key string) (float64, error) {
	val, err := r.GetString(ctx, key)
	if err != nil {
		return 0, err
	}

	result, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, errInvalidKeyType
	}

	return result, nil
}

// Delete removes a value by key
func (r *redisCacheManager) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return errCacheNotConnected
	}

	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in the cache
func (r *redisCacheManager) Exists(ctx context.Context, key string) (bool, error) {
	if r.client == nil {
		return false, errCacheNotConnected
	}

	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

// Keys returns all keys matching the given pattern
func (r *redisCacheManager) Keys(ctx context.Context, pattern string) ([]string, error) {
	if r.client == nil {
		return nil, errCacheNotConnected
	}

	return r.client.Keys(ctx, pattern).Result()
}

// Expire sets an expiration time for a key
func (r *redisCacheManager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if r.client == nil {
		return errCacheNotConnected
	}

	return r.client.Expire(ctx, key, expiration).Err()
}

// TTL returns the time to live for a key
func (r *redisCacheManager) TTL(ctx context.Context, key string) (time.Duration, error) {
	if r.client == nil {
		return 0, errCacheNotConnected
	}

	return r.client.TTL(ctx, key).Result()
}

// Clear removes all keys from the cache
func (r *redisCacheManager) Clear(ctx context.Context) error {
	if r.client == nil {
		return errCacheNotConnected
	}

	return r.client.FlushDB(ctx).Err()
}

// Ping checks if Redis is accessible
func (r *redisCacheManager) Ping(ctx context.Context) error {
	if r.client == nil {
		return errCacheNotConnected
	}

	return r.client.Ping(ctx).Err()
}

// SetMultiple stores multiple key-value pairs
func (r *redisCacheManager) SetMultiple(ctx context.Context, pairs map[string]interface{}, expiration time.Duration) error {
	if r.client == nil {
		return errCacheNotConnected
	}

	pipe := r.client.Pipeline()
	for key, value := range pairs {
		var data []byte
		var err error

		switch v := value.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			data, err = json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
			}
		}

		pipe.Set(ctx, key, data, expiration)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetMultiple retrieves multiple values by keys
func (r *redisCacheManager) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	if r.client == nil {
		return nil, errCacheNotConnected
	}

	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.StringCmd)

	for _, key := range keys {
		cmds[key] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for key, cmd := range cmds {
		val, err := cmd.Result()
		if err == nil {
			result[key] = val
		}
	}

	return result, nil
}

// DeleteMultiple removes multiple keys
func (r *redisCacheManager) DeleteMultiple(ctx context.Context, keys []string) error {
	if r.client == nil {
		return errCacheNotConnected
	}

	if len(keys) == 0 {
		return nil
	}

	return r.client.Del(ctx, keys...).Err()
}

// Increment increments a numeric value
func (r *redisCacheManager) Increment(ctx context.Context, key string, value int64) (int64, error) {
	if r.client == nil {
		return 0, errCacheNotConnected
	}

	return r.client.IncrBy(ctx, key, value).Result()
}

// Decrement decrements a numeric value
func (r *redisCacheManager) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	if r.client == nil {
		return 0, errCacheNotConnected
	}

	return r.client.DecrBy(ctx, key, value).Result()
}

// Close closes the Redis connection
func (r *redisCacheManager) Close() error {
	if r.client == nil {
		return nil
	}
	return r.client.Close()
}
