package ports

import (
	"context"
	"errors"

	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/common/genproto/stockpb"
	"github.com/mutition/go_start/stock/app"
	"github.com/sirupsen/logrus"
)

type GRPCServer struct {
	app app.Application
}

//new grpc server with application
func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

// CheckIfItemsInStock implements stockpb.StockServiceServer.
func (g GRPCServer) CheckIfItemsInStock(context.Context, *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	logrus.Info("CheckIfItemsInStock")
	defer func() {
		logrus.Info("CheckIfItemsInStock done")
	}()
	return nil, errors.New("not implemented")
}

// GetItems implements stockpb.StockServiceServer.
func (g GRPCServer) GetItems(context.Context, *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	logrus.Info("GetItems")
	defer func() {
		logrus.Info("GetItems done")
	}()
	fake := []*orderpb.Item{
		{
			Id: "i'm fake item",
		},
	}
	return &stockpb.GetItemsResponse{Items: fake}, nil
}

