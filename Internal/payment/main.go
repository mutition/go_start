package main

import (
	"context"
	"log"

	"github.com/mutition/go_start/common/broker"
	"github.com/mutition/go_start/common/config"
	"github.com/mutition/go_start/common/logging"
	"github.com/mutition/go_start/common/server"
	"github.com/mutition/go_start/payment/infrastructure/consumer"
	"github.com/mutition/go_start/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	//test stripe key
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logging.Init()
	logrus.Info(viper.GetString("stripe-key"))
	serverType := viper.GetString("payment.server-to-run")
	serviceName := viper.GetString("payment.service-name")
	paymentHandler := NewPaymentHandler()
	application, cleanup := service.NewApplication(ctx)
	defer cleanup()
	ch, closeCh := broker.ConnectToRabbitMQ(
		viper.GetString("rabbitmq.user"), viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"), viper.GetString("rabbitmq.port"))

	defer func() {
		_ = closeCh()
		_ = ch.Close()
	}()

	go consumer.NewConsumer(application).Listen(ch)

	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported server type:grpc")
	default:
		logrus.Panic("invalid server type")
	}
}
