package main

import (
	_ "common/config" // import for side effect to load configuration
	"common/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"order/ports"
)

func main() {
	serviceName := viper.GetString("order.service-name")
	server.RunHttpServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, NewHttpHandler(), ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil, // 中间件在 RunHttpServer 中统一设置
			ErrorHandler: nil,
		})
	})
}
