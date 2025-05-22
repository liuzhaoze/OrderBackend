package stock

import (
	"common/protobuf/stockpb"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGrpcClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	host, port := viper.GetString("stock.grpc-host"), viper.GetString("stock.grpc-port")
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.NewClient(address, dialOptions()...)
	if err != nil {
		return nil, nil, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func dialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}
