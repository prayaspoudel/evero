package config

import (
	"fmt"
	"testing"
)

// Example demonstrates how to use the config manager
func ExampleConfigManager() {
	// Method 1: Load config for a specific module
	configManager, err := LoadModuleConfig(ModuleHealth)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Access configuration values
	appName := configManager.GetString("app.name")
	webPort := configManager.GetInt("web.port")
	dbHost := configManager.GetString("database.host")

	fmt.Printf("App Name: %s\n", appName)
	fmt.Printf("Web Port: %d\n", webPort)
	fmt.Printf("DB Host: %s\n", dbHost)

	// Method 2: Manual configuration loading
	configManager2, err := NewConfigManagerFactory(InstanceViper)
	if err != nil {
		fmt.Printf("Error creating config manager: %v\n", err)
		return
	}

	if err := configManager2.Load(EnvLocal, ModuleHealth); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Check if a key exists
	if configManager2.IsSet("kafka.bootstrap.servers") {
		kafkaServers := configManager2.GetString("kafka.bootstrap.servers")
		fmt.Printf("Kafka Servers: %s\n", kafkaServers)
	}

	// Get all configuration as a map
	allConfig := configManager2.GetAll()
	fmt.Printf("Total config keys: %d\n", len(allConfig))
}

// Example_environmentVariables demonstrates environment variable support
func Example_environmentVariables() {
	configManager, err := NewConfigManagerFactory(InstanceViper)
	if err != nil {
		fmt.Printf("Error creating config manager: %v\n", err)
		return
	}

	// Cast to viperConfigManager to access additional methods
	if viperManager, ok := configManager.(*viperConfigManager); ok {
		// Enable environment variable support with prefix
		viperManager.SetEnvironmentVariables("APP")

		// Load configuration
		if err := viperManager.Load(EnvLocal, ModuleHealth); err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		// Now you can access config values that can be overridden by environment variables
		// For example, APP_WEB_PORT will override web.port
		webPort := viperManager.GetInt("web.port")
		fmt.Printf("Web Port (with env override): %d\n", webPort)
	}
}

// TestConfigManagerBasicFunctionality tests basic config manager functionality
func TestConfigManagerBasicFunctionality(t *testing.T) {
	// This is a basic test structure - you would need to set up proper test files
	configManager, err := NewConfigManagerFactory(InstanceViper)
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// Test that we can create an instance
	if configManager == nil {
		t.Fatal("Config manager should not be nil")
	}

	// Test environment validation
	if !ValidateEnvironment(EnvLocal) {
		t.Error("Local environment should be valid")
	}

	if ValidateEnvironment("invalid") {
		t.Error("Invalid environment should not be valid")
	}

	// Test module validation
	if !ValidateModule(ModuleHealth) {
		t.Error("Health module should be valid")
	}

	if ValidateModule("invalid") {
		t.Error("Invalid module should not be valid")
	}
}
