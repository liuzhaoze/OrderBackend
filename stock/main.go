package main

import (
	_ "common/config"
	"common/discovery"
	"common/protobuf/stockpb"
	"common/server"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := NewApplication(ctx)
	defer cleanup()

	deregisterConsul, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer deregisterConsul()

	server.RunGrpcServer(serviceName, func(s *grpc.Server) {
		stockpb.RegisterStockServiceServer(s, NewGrpcHandler(application))
	})
}
