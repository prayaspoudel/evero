package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type viperConfigManager struct {
	viper       *viper.Viper
	environment string
	module      string
}

// NewViperConfigManager creates a new Viper-based configuration manager
func NewViperConfigManager() (ConfigManager, error) {
	v := viper.New()
	return &viperConfigManager{
		viper: v,
	}, nil
}

// Load loads configuration from the specified environment and module
func (v *viperConfigManager) Load(environment, module string) error {
	v.environment = environment
	v.module = module

	// Reset viper instance
	v.viper = viper.New()

	// Set config type
	v.viper.SetConfigType("json")

	// Load base configuration from config/{module}/{environment}.json
	if module != "" {
		configPath := filepath.Join("config", module)
		v.viper.AddConfigPath(configPath)
		v.viper.SetConfigName(environment)

		if err := v.viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read base config for module %s and environment %s: %w", module, environment, err)
		}
	}

	// Load module-specific overrides from modules/{module}/config if they exist
	if module != "" {
		moduleConfigPath := filepath.Join("modules", module, "config")

		// Try to read module-specific config file
		moduleViper := viper.New()
		moduleViper.SetConfigType("json")
		moduleViper.AddConfigPath(moduleConfigPath)
		moduleViper.SetConfigName(environment)

		// If module-specific config exists, merge it
		if err := moduleViper.ReadInConfig(); err == nil {
			// Merge module-specific config into main config
			for key, value := range moduleViper.AllSettings() {
				v.viper.Set(key, value)
			}
		}
	}

	return nil
}

// Get retrieves a configuration value by key
func (v *viperConfigManager) Get(key string) interface{} {
	return v.viper.Get(key)
}

// GetString retrieves a string configuration value
func (v *viperConfigManager) GetString(key string) string {
	return v.viper.GetString(key)
}

// GetInt retrieves an integer configuration value
func (v *viperConfigManager) GetInt(key string) int {
	return v.viper.GetInt(key)
}

// GetBool retrieves a boolean configuration value
func (v *viperConfigManager) GetBool(key string) bool {
	return v.viper.GetBool(key)
}

// GetFloat64 retrieves a float64 configuration value
func (v *viperConfigManager) GetFloat64(key string) float64 {
	return v.viper.GetFloat64(key)
}

// IsSet checks if a configuration key is set
func (v *viperConfigManager) IsSet(key string) bool {
	return v.viper.IsSet(key)
}

// GetAll returns all configuration as a map
func (v *viperConfigManager) GetAll() map[string]interface{} {
	return v.viper.AllSettings()
}

// GetViper returns the underlying viper instance for advanced operations
func (v *viperConfigManager) GetViper() *viper.Viper {
	return v.viper
}

// LoadFromPaths loads configuration from specific paths (utility method)
func (v *viperConfigManager) LoadFromPaths(configName string, paths ...string) error {
	v.viper = viper.New()
	v.viper.SetConfigType("json")
	v.viper.SetConfigName(configName)

	for _, path := range paths {
		v.viper.AddConfigPath(path)
	}

	return v.viper.ReadInConfig()
}

// SetEnvironmentVariables enables environment variable support with prefix
func (v *viperConfigManager) SetEnvironmentVariables(prefix string) {
	v.viper.SetEnvPrefix(prefix)
	v.viper.AutomaticEnv()
	v.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// GetEnvironment returns the currently loaded environment
func (v *viperConfigManager) GetEnvironment() string {
	return v.environment
}

// GetModule returns the currently loaded module
func (v *viperConfigManager) GetModule() string {
	return v.module
}

// NewViper creates a new Viper configuration loader with simplified setup
// This provides backwards compatibility with the existing simple Viper usage
// configPath: path to config directory (e.g., "config/healthcare")
// configName: name of config file without extension (e.g., "local", "development", "production")
func NewViper(configPath, configName string) *viper.Viper {
	config := viper.New()

	config.SetConfigName(configName)
	config.SetConfigType("json")
	config.AddConfigPath(configPath)
	config.AddConfigPath("./../../" + configPath)
	config.AddConfigPath("./../" + configPath)
	config.AddConfigPath("./" + configPath)

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return config
}

// NewStructuredConfig creates a new structured configuration manager
// This is an alias for NewViperConfigManager for consistency with other infrastructure modules
func NewStructuredConfig() (ConfigManager, error) {
	return NewConfigManagerFactory(InstanceViper)
}
