# Infrastructure Setup Package

This package provides centralized infrastructure setup functions that can be reused across all modules (healthcare, insurance, etc.) in the evero project.

## Purpose

Previously, each module had its own infrastructure setup code (fiber.go, gorm.go, kafka.go, etc.) which led to code duplication. This package centralizes all infrastructure-related setup to promote:

- **Code Reusability**: Same infrastructure code can be used across healthcare, insurance, and other modules
- **Maintainability**: Changes to infrastructure setup only need to be made in one place
- **Consistency**: All modules use the same infrastructure setup patterns

## Available Functions

### Configuration
- **`NewViper(configPath, configName string) *viper.Viper`**
  - Loads configuration from JSON files
  - Parameters:
    - `configPath`: Path to config directory (e.g., "config/healthcare")
    - `configName`: Config file name without extension (e.g., "local", "development", "production")

### Logging
- **`NewLogger(viper *viper.Viper) *logrus.Logger`**
  - Creates a Logrus logger instance based on configuration
  - Reads log level from config: `log.level`

### Database
- **`NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB`**
  - Creates a PostgreSQL database connection using GORM
  - Reads database config from:
    - `database.username`
    - `database.password`
    - `database.host`
    - `database.port`
    - `database.name`
    - `database.pool.idle`
    - `database.pool.max`
    - `database.pool.lifetime`

### Validation
- **`NewValidator(viper *viper.Viper) *validator.Validate`**
  - Creates a go-playground validator instance

### Web Framework
- **`NewFiber(config *viper.Viper) *fiber.App`**
  - Creates a Fiber web application instance
  - Reads config from:
    - `app.name`: Application name
    - `web.prefork`: Enable/disable prefork mode

### Message Broker
- **`NewKafkaProducer(config *viper.Viper, log *logrus.Logger) sarama.SyncProducer`**
  - Creates a Kafka producer
  - Returns `nil` if `kafka.producer.enabled` is false
  - Reads config from:
    - `kafka.bootstrap.servers`
    - `kafka.producer.enabled`

- **`NewKafkaConsumerGroup(config *viper.Viper, log *logrus.Logger) sarama.ConsumerGroup`**
  - Creates a Kafka consumer group
  - Reads config from:
    - `kafka.bootstrap.servers`
    - `kafka.group.id`
    - `kafka.auto.offset.reset`

## Usage Example

### Healthcare Module Setup

```go
package setup

import (
    "fmt"
    
    infraSetup "github.com/prayaspoudel/infrastructure/setup"
    "github.com/prayaspoudel/modules/healthcare/config"
)

func Setup() {
    // Initialize infrastructure components
    viperConfig := infraSetup.NewViper("config/healthcare", "local")
    log := infraSetup.NewLogger(viperConfig)
    db := infraSetup.NewDatabase(viperConfig, log)
    validate := infraSetup.NewValidator(viperConfig)
    app := infraSetup.NewFiber(viperConfig)
    producer := infraSetup.NewKafkaProducer(viperConfig, log)

    // Bootstrap healthcare module with infrastructure
    config.Bootstrap(&config.BootstrapConfig{
        DB:       db,
        App:      app,
        Log:      log,
        Validate: validate,
        Config:   viperConfig,
        Producer: producer,
    })

    // Start server
    webPort := viperConfig.GetInt("web.port")
    err := app.Listen(fmt.Sprintf(":%d", webPort))
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

### Insurance Module Setup (Future)

The same infrastructure setup can be reused:

```go
package setup

import (
    infraSetup "github.com/prayaspoudel/infrastructure/setup"
    "github.com/prayaspoudel/modules/insurance/config"
)

func Setup() {
    // Same infrastructure setup code
    viperConfig := infraSetup.NewViper("config/insurance", "local")
    log := infraSetup.NewLogger(viperConfig)
    db := infraSetup.NewDatabase(viperConfig, log)
    // ... rest of the setup
    
    // Bootstrap insurance module
    config.Bootstrap(&config.BootstrapConfig{
        // ... insurance-specific bootstrap
    })
}
```

## Migration Notes

### Before
Each module had duplicate files:
- `modules/healthcare/config/fiber.go`
- `modules/healthcare/config/gorm.go`
- `modules/healthcare/config/kafka.go`
- `modules/healthcare/config/logrus.go`
- `modules/healthcare/config/validator.go`
- `modules/healthcare/config/viper.go`

### After
All infrastructure setup is centralized in:
- `infrastructure/setup/fiber.go`
- `infrastructure/setup/database.go`
- `infrastructure/setup/kafka.go`
- `infrastructure/setup/logger.go`
- `infrastructure/setup/validator.go`
- `infrastructure/setup/config.go`

Module-specific config only contains:
- `modules/healthcare/config/app.go` - Bootstrap function for healthcare module

## Benefits

1. **No Code Duplication**: Infrastructure code is written once and reused
2. **Easy to Maintain**: Update infrastructure in one place, all modules benefit
3. **Consistent Behavior**: All modules use the same infrastructure setup
4. **Easy to Add New Modules**: New modules can immediately use existing infrastructure
5. **Clear Separation**: Module-specific logic stays in modules, infrastructure logic is centralized
