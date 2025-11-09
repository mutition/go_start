package main

import (
	"context"
	"log"

	"github.com/mutition/go_start/common/config"
	"github.com/mutition/go_start/common/genproto/stockpb"
	"github.com/mutition/go_start/common/server"
	"github.com/mutition/go_start/stock/ports"
	"github.com/mutition/go_start/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application := service.NewApplication(ctx)

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
