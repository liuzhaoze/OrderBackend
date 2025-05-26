package order

import (
	_ "common/config"
	"common/discovery"
	"common/protobuf/orderpb"
	"context"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewOrderGrpcClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
	orderGrpcServerAddress, err := discovery.GetServiceAddress(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, nil, err
	}
	conn, err := grpc.NewClient(orderGrpcServerAddress, dialOptions()...)
	if err != nil {
		return nil, nil, err
	}
	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func dialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}
