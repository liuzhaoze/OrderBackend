package main

import (
	_ "common/config" // import for side effect to load configuration
	"common/server"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"order/ports"
)

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := NewApplication(ctx)
	defer cleanup()

	server.RunHttpServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, NewHttpHandler(application), ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil, // 中间件在 RunHttpServer 中统一设置
			ErrorHandler: nil,
		})
	})
}
