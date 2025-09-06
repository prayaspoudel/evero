package cache

import (
	"errors"
)

var (
	errInvalidCacheInstance = errors.New("invalid cache manager instance")
	errCacheNotConnected    = errors.New("cache manager not connected")
	errKeyNotFound          = errors.New("key not found in cache")
	errInvalidKeyType       = errors.New("invalid key type")
)

const (
	InstanceRedis int = iota
	InstanceInMemory
)

// NewCacheManagerFactory creates a new cache manager instance based on the specified type
func NewCacheManagerFactory(instance int, config *CacheConfig) (CacheManager, error) {
	switch instance {
	case InstanceRedis:
		return NewRedisCacheManager(config)
	case InstanceInMemory:
		return NewInMemoryCacheManager(config)
	default:
		return nil, errInvalidCacheInstance
	}
}
