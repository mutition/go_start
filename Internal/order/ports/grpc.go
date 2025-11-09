package ports

import (
	"context"

	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/order/app"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

// CreateOrder implements orderpb.OrderServiceServer.
func (g GRPCServer) CreateOrder(context.Context, *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// GetOrder implements orderpb.OrderServiceServer.
func (g GRPCServer) GetOrder(context.Context, *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	panic("unimplemented")
}

// UpdateOrder implements orderpb.OrderServiceServer.
func (g GRPCServer) UpdateOrder(context.Context, *orderpb.Order) (*emptypb.Empty, error) {
	panic("unimplemented")
}
