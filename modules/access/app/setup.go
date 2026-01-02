package access

import (
	"fmt"

	"github.com/prayaspoudel/infrastructure/config"
	"github.com/prayaspoudel/infrastructure/database"
	"github.com/prayaspoudel/infrastructure/logger"
	messagebroker "github.com/prayaspoudel/infrastructure/message-broker"
	"github.com/prayaspoudel/infrastructure/router"
	"github.com/prayaspoudel/infrastructure/validator"
)

func Setup() {
	// Initialize infrastructure components using specific infrastructure packages
	viperConfig := config.NewViper("config/access", "local")
	log := logger.NewLogger(viperConfig)
	db := database.NewDatabase(viperConfig, log)
	validate := validator.NewValidator(viperConfig)
	app := router.NewFiber(viperConfig)
	producer := messagebroker.NewKafkaProducer(viperConfig, log)

	// Bootstrap access module
	Bootstrap(&BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Producer: producer,
	})

	webPort := viperConfig.GetInt("web.port")
	if webPort == 0 {
		webPort = 8080 // Default port for SSO service
	}

	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
