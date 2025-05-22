package main

import (
	_ "common/config"
	"common/protobuf/stockpb"
	"common/server"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	server.RunGrpcServer(serviceName, func(s *grpc.Server) {
		stockpb.RegisterStockServiceServer(s, NewGrpcHandler())
	})
}
