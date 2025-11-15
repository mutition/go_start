package client

import (
	"context"
	"fmt"

	"github.com/mutition/go_start/common/discovery"
	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/common/genproto/stockpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	//服务发现
	grpcAddr, err := discovery.GetServiceGRPCAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, nil, err
	}
	if grpcAddr == "" {
		return nil, nil, fmt.Errorf("no grpc address found for service %s", viper.GetString("stock.service-name"))
	}
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

func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
	grpcAddr, err := discovery.GetServiceGRPCAddr(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, nil, err
	}
	if grpcAddr == "" {
		return nil, nil, fmt.Errorf("no grpc address found for service %s", viper.GetString("order.service-name"))
	}
	dialOptions, err := grpcDialOptions(grpcAddr)
	if err != nil {
		return nil, nil, err
	}
	conn, err := grpc.NewClient(grpcAddr, dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func grpcDialOptions(addr string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
}
