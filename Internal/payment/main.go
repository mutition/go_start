package main

import (
	"log"

	"github.com/mutition/go_start/common/config"
	"github.com/mutition/go_start/common/server"
	"github.com/mutition/go_start/common/logging"
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
	logrus.Info(viper.GetString("stripe-key"))
	serverType := viper.GetString("payment.server-to-run")
	serviceName := viper.GetString("payment.service-name")
	paymentHandler := NewPaymentHandler()
	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported server type:grpc")
	default:
		logrus.Panic("invalid server type")
	}
}
