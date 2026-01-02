package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prayaspoudel/infrastructure/config"
	"github.com/prayaspoudel/infrastructure/logger"
	messagebroker "github.com/prayaspoudel/infrastructure/message-broker"
	"github.com/prayaspoudel/modules/healthcare/delivery/messaging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viperConfig := config.NewViper("config/healthcare", "local")
	logger := logger.NewLogger(viperConfig)
	logger.Info("Starting worker service")

	ctx, cancel := context.WithCancel(context.Background())

	go RunUserConsumer(logger, viperConfig, ctx)
	go RunContactConsumer(logger, viperConfig, ctx)
	go RunAddressConsumer(logger, viperConfig, ctx)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	stop := false
	for !stop {
		select {
		case s := <-terminateSignals:
			logger.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
			cancel()
			stop = true
		}
	}

	time.Sleep(5 * time.Second) // wait for all consumers to finish processing
}

func RunAddressConsumer(logger *logrus.Logger, viperConfig *viper.Viper, ctx context.Context) {
	logger.Info("setup address consumer")
	addressConsumerGroup := messagebroker.NewKafkaConsumerGroup(viperConfig, logger)
	addressHandler := messaging.NewAddressConsumer(logger)
	messaging.ConsumeTopic(ctx, addressConsumerGroup, "addresses", logger, addressHandler.Consume)
}

func RunContactConsumer(logger *logrus.Logger, viperConfig *viper.Viper, ctx context.Context) {
	logger.Info("setup contact consumer")
	contactConsumerGroup := messagebroker.NewKafkaConsumerGroup(viperConfig, logger)
	contactHandler := messaging.NewContactConsumer(logger)
	messaging.ConsumeTopic(ctx, contactConsumerGroup, "contacts", logger, contactHandler.Consume)
}

func RunUserConsumer(logger *logrus.Logger, viperConfig *viper.Viper, ctx context.Context) {
	logger.Info("setup user consumer")
	userConsumerGroup := messagebroker.NewKafkaConsumerGroup(viperConfig, logger)
	userHandler := messaging.NewUserConsumer(logger)
	messaging.ConsumeTopic(ctx, userConsumerGroup, "users", logger, userHandler.Consume)
}
