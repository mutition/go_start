package ports

import (
	"context"

	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/order/app"
	"github.com/mutition/go_start/order/convertor"
	command "github.com/mutition/go_start/order/app/command"
	"github.com/mutition/go_start/order/app/query"
	domain "github.com/mutition/go_start/order/domain/order"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

// CreateOrder implements orderpb.OrderServiceServer.
func (g GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	_, err := g.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerId: request.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// GetOrder implements orderpb.OrderServiceServer.
func (g GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	o, err := g.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		OrderId:    request.OrderId,
		CustomerId: request.CustomerId,
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}
	return convertor.NewOrderConvertor().EntityToProto(o), nil
}

// UpdateOrder implements orderpb.OrderServiceServer.
func (g GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
	logrus.Infof("order_grpc || request_in || request: %v", request)
	order, err := domain.NewOrder(
		request.Id, request.CustomerId, request.Status, request.PaymentLink, 
		convertor.NewItemConvertor().ProtosToEntities(request.Items))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}
	_, err = g.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		},
	})
	return
}
