package ports

import (
	"context"

	"github.com/mutition/go_start/common/genproto/stockpb"
	"github.com/mutition/go_start/stock/app"
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
	panic("unimplemented")
}

// GetItems implements stockpb.StockServiceServer.
func (g GRPCServer) GetItems(context.Context, *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	panic("unimplemented")
}

