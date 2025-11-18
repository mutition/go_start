package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mutition/go_start/common/broker"
	"github.com/mutition/go_start/common/config"
	"github.com/mutition/go_start/common/discovery"
	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/common/logging"
	"github.com/mutition/go_start/common/server"
	"github.com/mutition/go_start/order/infrastructure/consumer"
	"github.com/mutition/go_start/order/ports"
	"github.com/mutition/go_start/order/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"github.com/mutition/go_start/common/tracing"
	// "github.com/mutition/go_start/order/ports"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err.Error())
	}
	// logrus.SetFormatter(&logrus.TextFormatter{
	//     DisableTimestamp: true,      // ❌ 不显示时间
	//     ForceColors:      true,      // ✅ 彩色输出
	//     FullTimestamp:    false,     // 简短格式
	//     PadLevelText:     true,      // 对齐 level
	//     TimestampFormat:  time.StampMilli,
	// })
}

func main() {
	logging.Init()
	serviceName := viper.GetString("order.service-name")
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
	ch, closeCh := broker.ConnectToRabbitMQ(
		viper.GetString("rabbitmq.user"), viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"), viper.GetString("rabbitmq.port"))

	defer func() {
		_ = closeCh()
		_ = ch.Close()
	}()
	deregisterfunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if err := deregisterfunc(); err != nil {
			logrus.Fatal(err)
		}
	}()
	application, cleanup := service.NewApplication(ctx)
	defer func() {
		_ = cleanup()
	}()

	go func() {
		_ = consumer.NewConsumer(application).Listen(ch)
	}()

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		orderpb.RegisterOrderServiceServer(server, ports.NewGRPCServer(application))
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		router.StaticFile("/payment/success", "../../public/success.html")
		ports.RegisterHandlersWithOptions(router, HTTPServer{app: application}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
