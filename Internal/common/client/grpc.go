package client

import (
	"context"
	"fmt"
	"time"

	"github.com/mutition/go_start/common/discovery"
	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	if !waitForStockGRPCCLient(10 * time.Second) {
		return nil, nil, fmt.Errorf("failed to connect to stock service")
	}
	grpcAddr, err := discovery.GetServiceGRPCAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, nil, err
	}
	if grpcAddr == "" {
		return nil, nil, fmt.Errorf("no grpc address found for service %s", viper.GetString("stock.service-name"))
	}
	dialOptions := grpcDialOptions(grpcAddr)
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
	if !waitForOrderGRPCCLient(10 * time.Second) {
		return nil, nil, fmt.Errorf("failed to connect to order service")
	}
	grpcAddr, err := discovery.GetServiceGRPCAddr(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, nil, err
	}
	if grpcAddr == "" {
		return nil, nil, fmt.Errorf("no grpc address found for service %s", viper.GetString("order.service-name"))
	}
	dialOptions := grpcDialOptions(grpcAddr)

	conn, err := grpc.NewClient(grpcAddr, dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func grpcDialOptions(_ string) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
}

func waitForOrderGRPCCLient(timeout time.Duration) bool {
	logrus.Info("waiting for order grpc client to be available")
	return waitForService(viper.GetString("order.service-name"), timeout)
}

func waitForStockGRPCCLient(timeout time.Duration) bool {
	logrus.Info("waiting for stock grpc client to be available")
	return waitForService(viper.GetString("stock.service-name"), timeout)
}


func waitForService( serviceName string, timeout time.Duration) bool {
	portAvailable := make(chan bool)
	timeoutCh := time.After(timeout)
	go func() {
		for {
			grpcAddr, err := discovery.GetServiceGRPCAddr(context.Background(), serviceName)
			if err == nil && grpcAddr != "" {
				portAvailable <- true
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	select {
		case <-portAvailable:
			return true
		case <-timeoutCh:
			return false
	}
}