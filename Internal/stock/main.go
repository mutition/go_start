package main

import (
	"context"
	"log"

	"github.com/mutition/go_start/common/config"
	"github.com/mutition/go_start/common/discovery"
	"github.com/mutition/go_start/common/genproto/stockpb"
	"github.com/mutition/go_start/common/logging"
	"github.com/mutition/go_start/common/server"
	"github.com/mutition/go_start/stock/ports"
	"github.com/mutition/go_start/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"github.com/mutition/go_start/common/tracing"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	logging.Init()
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	shutdownJaeger, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if err := shutdownJaeger(ctx); err != nil {
			logrus.Fatal(err)
		}
	}()
	application := service.NewApplication(ctx)

	deregisterfunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if err := deregisterfunc(); err != nil {
			logrus.Fatal(err)
		}
	}()

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			stockpb.RegisterStockServiceServer(server, ports.NewGRPCServer(application))
		})
	case "http":

	default:
		logrus.Panicf("invalid server type: %s", serverType)
	}
}
