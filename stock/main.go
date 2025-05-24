package main

import (
	_ "common/config"
	"common/protobuf/stockpb"
	"common/server"
	"context"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := NewApplication(ctx)
	defer cleanup()

	server.RunGrpcServer(serviceName, func(s *grpc.Server) {
		stockpb.RegisterStockServiceServer(s, NewGrpcHandler(application))
	})
}
