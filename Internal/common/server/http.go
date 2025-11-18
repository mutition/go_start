package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).Get("http-addr")
	RunHTTPServerOnaddr(addr, wrapper, serviceName)
}

func RunHTTPServerOnaddr(addr any, wrapper func(router *gin.Engine), serviceName string) {
	apiRouter := gin.New()
	setMiddleware(apiRouter, serviceName)
	wrapper(apiRouter)
	apiRouter.Group("/api")
	apiRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})
	if err := apiRouter.Run(addr.(string)); err != nil {
		panic(err)
	}
}

func setMiddleware(router *gin.Engine, serviceName string) {
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(serviceName))
}
