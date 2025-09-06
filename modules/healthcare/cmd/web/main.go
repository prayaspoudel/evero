package main

import (
	"os"
	"time"

	"github.com/prayaspoudel/infrastructure/database"
	"github.com/prayaspoudel/infrastructure/logger"
	"github.com/prayaspoudel/infrastructure/router"
	"github.com/prayaspoudel/infrastructure/validator"
)

func main() {
	// viperConfig := config.NewViper()
	// log := config.NewLogger(viperConfig)
	// db := config.NewDatabase(viperConfig, log)
	// validate := config.NewValidator(viperConfig)
	// app := config.NewFiber(viperConfig)
	// producer := config.NewKafkaProducer(viperConfig, log)

	// config.Bootstrap(&config.BootstrapConfig{
	// 	DB:       db,
	// 	App:      app,
	// 	Log:      log,
	// 	Validate: validate,
	// 	Config:   viperConfig,
	// 	Producer: producer,
	// })

	// webPort := viperConfig.GetInt("web.port")
	// err := app.Listen(fmt.Sprintf(":%d", webPort))
	// if err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
	Run()
}

func Run() {
	var app = router.NewConfig().
		Name(os.Getenv("APP_NAME")).
		ContextTimeout(10 * time.Second).
		Logger(logger.InstanceLogrusLogger).
		Validator(validator.InstanceGoPlayground).
		DBSQL(database.InstancePostgres).
		DBNoSQL(database.InstanceMongoDB)

	app.WebServerPort(os.Getenv("APP_PORT")).
		WebServer(router.InstanceGorillaMux).
		Start()
}
