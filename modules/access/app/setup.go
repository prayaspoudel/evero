package access

import (
	"fmt"

	infraSetup "github.com/prayaspoudel/infrastructure/setup"
)

func Setup() {
	// Initialize infrastructure components using centralized setup
	viperConfig := infraSetup.NewViper("config/access", "local")
	log := infraSetup.NewLogger(viperConfig)
	db := infraSetup.NewDatabase(viperConfig, log)
	validate := infraSetup.NewValidator(viperConfig)
	app := infraSetup.NewFiber(viperConfig)
	producer := infraSetup.NewKafkaProducer(viperConfig, log)

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
