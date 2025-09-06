package cache_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/prayaspoudel/infrastructure/cache"
)

func TestInMemoryCacheManager(t *testing.T) {
	config := &cache.CacheConfig{
		DefaultExpiration: 5 * time.Minute,
		CleanupInterval:   10 * time.Minute,
		MaxSize:           100,
	}

	cacheManager, err := cache.NewCacheManagerFactory(cache.InstanceInMemory, config)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()
	err = cacheManager.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer cacheManager.Close()

	// Test Set and Get
	key := "test-key"
	value := "test-value"
	expiration := time.Minute

	err = cacheManager.Set(ctx, key, value, expiration)
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	retrievedValue, err := cacheManager.GetString(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}

	if retrievedValue != value {
		t.Errorf("Expected %s, got %s", value, retrievedValue)
	}

	// Test Exists
	exists, err := cacheManager.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}

	if !exists {
		t.Error("Key should exist")
	}

	// Test Delete
	err = cacheManager.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	exists, err = cacheManager.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Failed to check existence after delete: %v", err)
	}

	if exists {
		t.Error("Key should not exist after delete")
	}
}

func TestCacheManagerTypes(t *testing.T) {
	config := &cache.CacheConfig{
		DefaultExpiration: time.Hour,
		MaxSize:           100,
	}

	cacheManager, err := cache.NewCacheManagerFactory(cache.InstanceInMemory, config)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()
	err = cacheManager.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer cacheManager.Close()

	// Test different data types
	tests := []struct {
		key    string
		value  interface{}
		getter func(string) (interface{}, error)
	}{
		{"string", "hello", func(k string) (interface{}, error) { return cacheManager.GetString(ctx, k) }},
		{"int", 42, func(k string) (interface{}, error) { return cacheManager.GetInt(ctx, k) }},
		{"bool", true, func(k string) (interface{}, error) { return cacheManager.GetBool(ctx, k) }},
		{"float", 3.14, func(k string) (interface{}, error) { return cacheManager.GetFloat64(ctx, k) }},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			err := cacheManager.Set(ctx, test.key, test.value, time.Hour)
			if err != nil {
				t.Fatalf("Failed to set %s: %v", test.key, err)
			}

			retrieved, err := test.getter(test.key)
			if err != nil {
				t.Fatalf("Failed to get %s: %v", test.key, err)
			}

			// For comparison, convert both to string if needed
			if retrieved != test.value {
				t.Errorf("Expected %v, got %v", test.value, retrieved)
			}
		})
	}
}

func TestCacheManagerBatchOperations(t *testing.T) {
	config := &cache.CacheConfig{
		DefaultExpiration: time.Hour,
		MaxSize:           100,
	}

	cacheManager, err := cache.NewCacheManagerFactory(cache.InstanceInMemory, config)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()
	err = cacheManager.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer cacheManager.Close()

	// Test SetMultiple
	pairs := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	err = cacheManager.SetMultiple(ctx, pairs, time.Hour)
	if err != nil {
		t.Fatalf("Failed to set multiple: %v", err)
	}

	// Test GetMultiple
	keys := []string{"key1", "key2", "key3"}
	values, err := cacheManager.GetMultiple(ctx, keys)
	if err != nil {
		t.Fatalf("Failed to get multiple: %v", err)
	}

	if len(values) != len(pairs) {
		t.Errorf("Expected %d values, got %d", len(pairs), len(values))
	}

	for key, expectedValue := range pairs {
		if value, exists := values[key]; !exists {
			t.Errorf("Key %s not found in results", key)
		} else if value != expectedValue {
			t.Errorf("Key %s: expected %v, got %v", key, expectedValue, value)
		}
	}

	// Test DeleteMultiple
	err = cacheManager.DeleteMultiple(ctx, keys)
	if err != nil {
		t.Fatalf("Failed to delete multiple: %v", err)
	}

	// Verify deletion
	for _, key := range keys {
		exists, err := cacheManager.Exists(ctx, key)
		if err != nil {
			t.Fatalf("Failed to check existence for %s: %v", key, err)
		}
		if exists {
			t.Errorf("Key %s should not exist after delete", key)
		}
	}
}

func TestCacheManagerIncrement(t *testing.T) {
	config := &cache.CacheConfig{
		DefaultExpiration: time.Hour,
		MaxSize:           100,
	}

	cacheManager, err := cache.NewCacheManagerFactory(cache.InstanceInMemory, config)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()
	err = cacheManager.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer cacheManager.Close()

	key := "counter"

	// Test increment on non-existent key
	newValue, err := cacheManager.Increment(ctx, key, 1)
	if err != nil {
		t.Fatalf("Failed to increment: %v", err)
	}

	if newValue != 1 {
		t.Errorf("Expected 1, got %d", newValue)
	}

	// Test increment on existing key
	newValue, err = cacheManager.Increment(ctx, key, 5)
	if err != nil {
		t.Fatalf("Failed to increment: %v", err)
	}

	if newValue != 6 {
		t.Errorf("Expected 6, got %d", newValue)
	}

	// Test decrement
	newValue, err = cacheManager.Decrement(ctx, key, 2)
	if err != nil {
		t.Fatalf("Failed to decrement: %v", err)
	}

	if newValue != 4 {
		t.Errorf("Expected 4, got %d", newValue)
	}
}

func BenchmarkInMemoryCache(b *testing.B) {
	config := &cache.CacheConfig{
		DefaultExpiration: time.Hour,
		MaxSize:           10000,
	}

	cacheManager, err := cache.NewCacheManagerFactory(cache.InstanceInMemory, config)
	if err != nil {
		b.Fatalf("Failed to create cache manager: %v", err)
	}

	ctx := context.Background()
	err = cacheManager.Connect(ctx)
	if err != nil {
		b.Fatalf("Failed to connect: %v", err)
	}
	defer cacheManager.Close()

	b.ResetTimer()

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			cacheManager.Set(ctx, key, value, time.Hour)
		}
	})

	b.Run("Get", func(b *testing.B) {
		// Pre-populate cache
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			cacheManager.Set(ctx, key, value, time.Hour)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key-%d", i%1000)
			cacheManager.GetString(ctx, key)
		}
	})
}
