package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mutition/go_start/common/config"
	"github.com/mutition/go_start/common/server"
	"github.com/mutition/go_start/order/ports"

	// "github.com/mutition/go_start/order/ports"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{},ports.GinServerOptions{
			BaseURL: "/api",
			Middlewares: nil,
			ErrorHandler: nil,
		})
	})
}
