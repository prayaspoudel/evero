package cache

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration int64
}

type inMemoryCacheManager struct {
	items           map[string]*cacheItem
	mutex           sync.RWMutex
	config          *CacheConfig
	cleanupInterval time.Duration
	stopCleanup     chan bool
}

// NewInMemoryCacheManager creates a new in-memory cache manager
func NewInMemoryCacheManager(config *CacheConfig) (CacheManager, error) {
	if config == nil {
		config = &CacheConfig{
			DefaultExpiration: 5 * time.Minute,
			CleanupInterval:   10 * time.Minute,
			MaxSize:           1000,
		}
	}

	// Set defaults if not specified
	if config.DefaultExpiration == 0 {
		config.DefaultExpiration = 5 * time.Minute
	}
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 10 * time.Minute
	}
	if config.MaxSize == 0 {
		config.MaxSize = 1000
	}

	manager := &inMemoryCacheManager{
		items:           make(map[string]*cacheItem),
		config:          config,
		cleanupInterval: config.CleanupInterval,
		stopCleanup:     make(chan bool),
	}

	return manager, nil
}

// Connect starts the cleanup goroutine
func (m *inMemoryCacheManager) Connect(ctx context.Context) error {
	go m.startCleanup()
	return nil
}

// Disconnect stops the cleanup goroutine
func (m *inMemoryCacheManager) Disconnect(ctx context.Context) error {
	close(m.stopCleanup)
	return nil
}

// isExpired checks if an item has expired
func (item *cacheItem) isExpired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.expiration
}

// startCleanup starts the cleanup goroutine
func (m *inMemoryCacheManager) startCleanup() {
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanup()
		case <-m.stopCleanup:
			return
		}
	}
}

// cleanup removes expired items
func (m *inMemoryCacheManager) cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for key, item := range m.items {
		if item.isExpired() {
			delete(m.items, key)
		}
	}
}

// Set stores a value with the given key and expiration time
func (m *inMemoryCacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if we need to make room
	if len(m.items) >= m.config.MaxSize {
		// Simple eviction: remove first expired item found, or oldest item
		for k, item := range m.items {
			if item.isExpired() {
				delete(m.items, k)
				break
			}
		}
		// If still at max capacity, remove one item (simple FIFO)
		if len(m.items) >= m.config.MaxSize {
			for k := range m.items {
				delete(m.items, k)
				break
			}
		}
	}

	var exp int64
	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	}

	m.items[key] = &cacheItem{
		value:      value,
		expiration: exp,
	}

	return nil
}

// Get retrieves a value by key
func (m *inMemoryCacheManager) Get(ctx context.Context, key string) (interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, found := m.items[key]
	if !found {
		return nil, errKeyNotFound
	}

	if item.isExpired() {
		m.mutex.RUnlock()
		m.mutex.Lock()
		delete(m.items, key)
		m.mutex.Unlock()
		m.mutex.RLock()
		return nil, errKeyNotFound
	}

	return item.value, nil
}

// GetString retrieves a string value by key
func (m *inMemoryCacheManager) GetString(ctx context.Context, key string) (string, error) {
	value, err := m.Get(ctx, key)
	if err != nil {
		return "", err
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// GetInt retrieves an integer value by key
func (m *inMemoryCacheManager) GetInt(ctx context.Context, key string) (int, error) {
	value, err := m.Get(ctx, key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, errInvalidKeyType
	}
}

// GetBool retrieves a boolean value by key
func (m *inMemoryCacheManager) GetBool(ctx context.Context, key string) (bool, error) {
	value, err := m.Get(ctx, key)
	if err != nil {
		return false, err
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, errInvalidKeyType
	}
}

// GetFloat64 retrieves a float64 value by key
func (m *inMemoryCacheManager) GetFloat64(ctx context.Context, key string) (float64, error) {
	value, err := m.Get(ctx, key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, errInvalidKeyType
	}
}

// Delete removes a value by key
func (m *inMemoryCacheManager) Delete(ctx context.Context, key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.items, key)
	return nil
}

// Exists checks if a key exists in the cache
func (m *inMemoryCacheManager) Exists(ctx context.Context, key string) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, found := m.items[key]
	if !found {
		return false, nil
	}

	if item.isExpired() {
		m.mutex.RUnlock()
		m.mutex.Lock()
		delete(m.items, key)
		m.mutex.Unlock()
		m.mutex.RLock()
		return false, nil
	}

	return true, nil
}

// Keys returns all keys matching the given pattern (simple prefix matching for in-memory)
func (m *inMemoryCacheManager) Keys(ctx context.Context, pattern string) ([]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var keys []string
	for key, item := range m.items {
		if !item.isExpired() {
			// Simple pattern matching - just check if key contains pattern
			// For more sophisticated pattern matching, you'd need a proper glob library
			if pattern == "*" || key == pattern {
				keys = append(keys, key)
			}
		}
	}

	return keys, nil
}

// Expire sets an expiration time for a key
func (m *inMemoryCacheManager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	item, found := m.items[key]
	if !found {
		return errKeyNotFound
	}

	if item.isExpired() {
		delete(m.items, key)
		return errKeyNotFound
	}

	item.expiration = time.Now().Add(expiration).UnixNano()
	return nil
}

// TTL returns the time to live for a key
func (m *inMemoryCacheManager) TTL(ctx context.Context, key string) (time.Duration, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, found := m.items[key]
	if !found {
		return 0, errKeyNotFound
	}

	if item.isExpired() {
		return 0, errKeyNotFound
	}

	if item.expiration == 0 {
		return -1, nil // No expiration
	}

	ttl := time.Duration(item.expiration - time.Now().UnixNano())
	if ttl < 0 {
		return 0, nil
	}

	return ttl, nil
}

// Clear removes all keys from the cache
func (m *inMemoryCacheManager) Clear(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.items = make(map[string]*cacheItem)
	return nil
}

// Ping always returns nil for in-memory cache
func (m *inMemoryCacheManager) Ping(ctx context.Context) error {
	return nil
}

// SetMultiple stores multiple key-value pairs
func (m *inMemoryCacheManager) SetMultiple(ctx context.Context, pairs map[string]interface{}, expiration time.Duration) error {
	for key, value := range pairs {
		if err := m.Set(ctx, key, value, expiration); err != nil {
			return err
		}
	}
	return nil
}

// GetMultiple retrieves multiple values by keys
func (m *inMemoryCacheManager) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, key := range keys {
		value, err := m.Get(ctx, key)
		if err == nil {
			result[key] = value
		}
	}
	return result, nil
}

// DeleteMultiple removes multiple keys
func (m *inMemoryCacheManager) DeleteMultiple(ctx context.Context, keys []string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, key := range keys {
		delete(m.items, key)
	}
	return nil
}

// Increment increments a numeric value
func (m *inMemoryCacheManager) Increment(ctx context.Context, key string, value int64) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	item, found := m.items[key]
	if !found {
		// Create new item with the increment value
		m.items[key] = &cacheItem{
			value:      value,
			expiration: 0,
		}
		return value, nil
	}

	if item.isExpired() {
		delete(m.items, key)
		m.items[key] = &cacheItem{
			value:      value,
			expiration: 0,
		}
		return value, nil
	}

	// Try to convert existing value to int64
	switch v := item.value.(type) {
	case int64:
		newValue := v + value
		item.value = newValue
		return newValue, nil
	case int:
		newValue := int64(v) + value
		item.value = newValue
		return newValue, nil
	case string:
		current, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errInvalidKeyType
		}
		newValue := current + value
		item.value = newValue
		return newValue, nil
	default:
		return 0, errInvalidKeyType
	}
}

// Decrement decrements a numeric value
func (m *inMemoryCacheManager) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return m.Increment(ctx, key, -value)
}

// Close closes the in-memory cache manager
func (m *inMemoryCacheManager) Close() error {
	close(m.stopCleanup)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.items = nil
	return nil
}
