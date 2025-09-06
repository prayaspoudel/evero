// Package config provides configuration management functionality for loading
// and accessing application configuration from various sources including
// environment-specific files and module-specific configurations.
package config

import "github.com/spf13/viper"

// ConfigManager interface defines the contract for configuration management
type ConfigManager interface {
	// Load loads configuration from the specified environment and module
	Load(environment, module string) error

	// Get retrieves a configuration value by key
	Get(key string) interface{}

	// GetString retrieves a string configuration value
	GetString(key string) string

	// GetInt retrieves an integer configuration value
	GetInt(key string) int

	// GetBool retrieves a boolean configuration value
	GetBool(key string) bool

	// GetFloat64 retrieves a float64 configuration value
	GetFloat64(key string) float64

	// IsSet checks if a configuration key is set
	IsSet(key string) bool

	// GetAll returns all configuration as a map
	GetAll() map[string]interface{}

	// GetViper returns the underlying viper instance for advanced operations
	GetViper() *viper.Viper
}
