package server

import (
	"common/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RunHttpServer(serviceName string, register func(router *gin.Engine)) {
	host := viper.Sub(serviceName).GetString("http-host")
	if host == "" {
		logrus.Panicln("Empty http host")
	}
	port := viper.Sub(serviceName).GetString("http-port")
	if port == "" {
		logrus.Panicln("Empty http port")
	}
	address := fmt.Sprintf("%s:%s", host, port)

	router := gin.New()
	// 设置中间件必须在注册路由之前
	setMiddlewares(router, serviceName)
	register(router)
	if err := router.Run(address); err != nil {
		logrus.Panicln(err)
	}
}

func setMiddlewares(router *gin.Engine, serviceName string) {
	router.Use(middleware.LogHttpRequest(logrus.NewEntry(logrus.StandardLogger())))
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(serviceName))
}
