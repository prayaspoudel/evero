package test

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/prayaspoudel/infrastructure/config"
	"github.com/prayaspoudel/infrastructure/database"
	"github.com/prayaspoudel/infrastructure/logger"
	messagebroker "github.com/prayaspoudel/infrastructure/message-broker"
	"github.com/prayaspoudel/infrastructure/router"
	infraValidator "github.com/prayaspoudel/infrastructure/validator"
	healthcare "github.com/prayaspoudel/modules/healthcare/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var app *fiber.App

var db *gorm.DB

var viperConfig *viper.Viper

var log *logrus.Logger

var validate *validator.Validate

func init() {
	viperConfig = config.NewViper("config/healthcare", "local")
	log = logger.NewLogger(viperConfig)
	validate = infraValidator.NewValidator(viperConfig)
	app = router.NewFiber(viperConfig)
	db = database.NewDatabase(viperConfig, log)
	producer := messagebroker.NewKafkaProducer(viperConfig, log)

	healthcare.Bootstrap(&healthcare.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Producer: producer,
	})
}
