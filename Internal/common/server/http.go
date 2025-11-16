package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).Get("http-addr")
	RunHTTPServerOnaddr(addr, wrapper)
	router := gin.Default()
	wrapper(router)
	err := router.Run(addr.(string))
	if err != nil {
		panic(err)
	}
}

func RunHTTPServerOnaddr(addr any, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	wrapper(apiRouter)
	apiRouter.Group("/api")
	apiRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})
	if err := apiRouter.Run(addr.(string)); err != nil {
		panic(err)
	}
}
