package stock

import (
	_ "common/config"
	"common/discovery"
	"common/protobuf/stockpb"
	"context"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGrpcClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	stockGrpcServerAddress, err := discovery.GetServiceAddress(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, nil, err
	}
	conn, err := grpc.NewClient(stockGrpcServerAddress, dialOptions()...)
	if err != nil {
		return nil, nil, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func dialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
}
