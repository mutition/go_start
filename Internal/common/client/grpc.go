package client

import (
	"context"

	"github.com/mutition/go_start/common/genproto/stockpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	grpcAddr := viper.GetString("stock.grpc-addr")
	dialOptions, err := grpcDialOptions(grpcAddr)
	if err != nil {
		return nil, nil, err
	}
	conn, err := grpc.NewClient(grpcAddr, dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil

}

func grpcDialOptions(addr string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
}