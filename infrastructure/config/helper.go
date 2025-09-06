package config

import (
	"fmt"
	"os"
)

// Environment constants
const (
	EnvLocal       = "local"
	EnvDevelopment = "development"
	EnvStage       = "stage"
	EnvProduction  = "production"
)

// Module constants
const (
	ModuleHealth    = "healthcare"
	ModuleInsurance = "insurance"
)

// GetEnvironment gets environment from environment variable or defaults to local
func GetEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = os.Getenv("ENVIRONMENT")
	}
	if env == "" {
		return EnvLocal
	}
	return env
}

// GetModule gets module from environment variable
func GetModule() string {
	return os.Getenv("MODULE")
}

// LoadModuleConfig is a convenience function to load configuration for a specific module
func LoadModuleConfig(module string) (ConfigManager, error) {
	environment := GetEnvironment()

	configManager, err := NewConfigManagerFactory(InstanceViper)
	if err != nil {
		return nil, fmt.Errorf("failed to create config manager: %w", err)
	}

	if err := configManager.Load(environment, module); err != nil {
		return nil, fmt.Errorf("failed to load config for module %s in environment %s: %w", module, environment, err)
	}

	return configManager, nil
}

// LoadConfig is a convenience function to load configuration with auto-detection
func LoadConfig() (ConfigManager, error) {
	environment := GetEnvironment()
	module := GetModule()

	if module == "" {
		return nil, fmt.Errorf("module not specified in MODULE environment variable")
	}

	configManager, err := NewConfigManagerFactory(InstanceViper)
	if err != nil {
		return nil, fmt.Errorf("failed to create config manager: %w", err)
	}

	if err := configManager.Load(environment, module); err != nil {
		return nil, fmt.Errorf("failed to load config for module %s in environment %s: %w", module, environment, err)
	}

	return configManager, nil
}

// ValidateEnvironment checks if the provided environment is valid
func ValidateEnvironment(env string) bool {
	switch env {
	case EnvLocal, EnvDevelopment, EnvStage, EnvProduction:
		return true
	default:
		return false
	}
}

// ValidateModule checks if the provided module is valid
func ValidateModule(module string) bool {
	switch module {
	case ModuleHealth, ModuleInsurance:
		return true
	default:
		return false
	}
}
