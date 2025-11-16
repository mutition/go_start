package ports

import (
	"context"

	"github.com/mutition/go_start/common/genproto/stockpb"
	"github.com/mutition/go_start/stock/app"
	"github.com/mutition/go_start/stock/app/query"
)

type GRPCServer struct {
	app app.Application
}

// new grpc server with application
func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

// GetItems implements stockpb.StockServiceServer.
func (g GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	items, err := g.app.Queries.GetItems.Handle(ctx, query.GetItems{ItemIds: request.ItemIds})
	if err != nil {
		return nil, err
	}
	return &stockpb.GetItemsResponse{Items: items}, nil
}

// CheckIfItemsInStock implements stockpb.StockServiceServer.
func (g GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	items, err := g.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: request.Items})
	if err != nil {
		return nil, err
	}
	return &stockpb.CheckIfItemsInStockResponse{Items: items, InStock: 1}, nil
}
