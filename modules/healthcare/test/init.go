package test

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	infraSetup "github.com/prayaspoudel/infrastructure/setup"
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
	viperConfig = infraSetup.NewViper("config/healthcare", "local")
	log = infraSetup.NewLogger(viperConfig)
	validate = infraSetup.NewValidator(viperConfig)
	app = infraSetup.NewFiber(viperConfig)
	db = infraSetup.NewDatabase(viperConfig, log)
	producer := infraSetup.NewKafkaProducer(viperConfig, log)

	healthcare.Bootstrap(&healthcare.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Producer: producer,
	})
}
