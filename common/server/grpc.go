package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
)

func RunGrpcServer(serviceName string, register func(server *grpc.Server)) {
	host := viper.Sub(serviceName).GetString("grpc-host")
	if host == "" {
		logrus.Panicln("Empty grpc host")
	}
	port := viper.Sub(serviceName).GetString("grpc-port")
	if port == "" {
		logrus.Panicln("Empty grpc port")
	}
	address := fmt.Sprintf("%s:%s", host, port)

	// TODO: gRPC options
	server := grpc.NewServer()
	register(server)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logrus.Panicln(err)
	}
	if err = server.Serve(listener); err != nil {
		logrus.Panicln(err)
	}
}
