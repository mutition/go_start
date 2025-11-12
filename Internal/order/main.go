package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mutition/go_start/common/config"
	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/common/server"
	"github.com/mutition/go_start/order/ports"
	"github.com/mutition/go_start/order/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"github.com/mutition/go_start/common/discovery"

	// "github.com/mutition/go_start/order/ports"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err.Error())
	}
	logrus.SetFormatter(&logrus.TextFormatter{
        DisableTimestamp: true,      // ❌ 不显示时间
        ForceColors:      true,      // ✅ 彩色输出
        FullTimestamp:    false,     // 简短格式
        PadLevelText:     true,      // 对齐 level
        TimestampFormat:  time.StampMilli,
    })
}

func main() {
	serviceName := viper.GetString("order.service-name")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	defer cleanup()

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		orderpb.RegisterOrderServiceServer(server, ports.NewGRPCServer(application))
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{app: application}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
