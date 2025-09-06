# Config Manager

The config manager provides a centralized way to load and manage application configuration from various sources including environment-specific files and module-specific configurations.

## Features

- Load configuration from environment-specific JSON files
- Support for module-specific configuration overrides
- Environment variable support with prefix
- Type-safe configuration access methods
- Factory pattern for easy instantiation
- Helper functions for common usage patterns

## Directory Structure

The config manager expects the following directory structure:

```
config/
├── health/
│   ├── local.json
│   ├── development.json
│   ├── stage.json
│   └── production.json
└── insurance/
    ├── local.json
    ├── development.json
    ├── stage.json
    └── production.json

modules/
├── health/
│   └── config/
│       ├── local.json (optional overrides)
│       ├── development.json (optional overrides)
│       ├── stage.json (optional overrides)
│       └── production.json (optional overrides)
└── insurance/
    └── config/
        ├── local.json (optional overrides)
        ├── development.json (optional overrides)
        ├── stage.json (optional overrides)
        └── production.json (optional overrides)
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "your-project/infrastructure/config"
)

func main() {
    // Load configuration for health module
    configManager, err := config.LoadModuleConfig(config.ModuleHealth)
    if err != nil {
        log.Fatal(err)
    }
    
    // Access configuration values
    appName := configManager.GetString("app.name")
    webPort := configManager.GetInt("web.port")
    dbHost := configManager.GetString("database.host")
    
    fmt.Printf("App: %s, Port: %d, DB: %s\n", appName, webPort, dbHost)
}
```

### Manual Configuration Loading

```go
package main

import (
    "log"
    
    "your-project/infrastructure/config"
)

func main() {
    // Create config manager instance
    configManager, err := config.NewConfigManagerFactory(config.InstanceViper)
    if err != nil {
        log.Fatal(err)
    }
    
    // Load configuration for specific environment and module
    err = configManager.Load(config.EnvLocal, config.ModuleHealth)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use configuration
    if configManager.IsSet("database.host") {
        dbHost := configManager.GetString("database.host")
        fmt.Printf("Database Host: %s\n", dbHost)
    }
}
```

### Environment Variable Support

```go
package main

import (
    "log"
    
    "your-project/infrastructure/config"
)

func main() {
    configManager, err := config.NewConfigManagerFactory(config.InstanceViper)
    if err != nil {
        log.Fatal(err)
    }
    
    // Cast to access additional methods
    if viperManager, ok := configManager.(*config.viperConfigManager); ok {
        // Enable environment variable support with prefix
        viperManager.SetEnvironmentVariables("APP")
        
        // Load configuration
        err = viperManager.Load(config.EnvLocal, config.ModuleHealth)
        if err != nil {
            log.Fatal(err)
        }
        
        // Now APP_WEB_PORT environment variable can override web.port
        webPort := viperManager.GetInt("web.port")
        fmt.Printf("Web Port: %d\n", webPort)
    }
}
```

### Auto-Detection from Environment Variables

```go
package main

import (
    "log"
    "os"
    
    "your-project/infrastructure/config"
)

func main() {
    // Set environment variables
    os.Setenv("APP_ENV", "development")
    os.Setenv("MODULE", "health")
    
    // Load configuration with auto-detection
    configManager, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    // Use configuration
    appName := configManager.GetString("app.name")
    fmt.Printf("App Name: %s\n", appName)
}
```

## Environment Variables

- `APP_ENV` or `ENVIRONMENT`: Sets the environment (local, development, stage, production)
- `MODULE`: Sets the module name (health, insurance)
- `APP_*`: Configuration overrides when environment variable support is enabled

## Configuration Methods

### Available Methods

- `Load(environment, module string) error`: Load configuration for specific environment and module
- `Get(key string) interface{}`: Get any configuration value
- `GetString(key string) string`: Get string value
- `GetInt(key string) int`: Get integer value
- `GetBool(key string) bool`: Get boolean value
- `GetFloat64(key string) float64`: Get float64 value
- `IsSet(key string) bool`: Check if key exists
- `GetAll() map[string]interface{}`: Get all configuration as map
- `GetViper() *viper.Viper`: Get underlying viper instance

## Constants

### Environments
- `config.EnvLocal`
- `config.EnvDevelopment`
- `config.EnvStage`
- `config.EnvProduction`

### Modules
- `config.ModuleHealth`
- `config.ModuleInsurance`

## Configuration Loading Priority

1. Base configuration from `config/{module}/{environment}.json`
2. Module-specific overrides from `modules/{module}/config/{environment}.json` (if exists)
3. Environment variables (if enabled with prefix)

Later sources override earlier ones for the same configuration keys.
