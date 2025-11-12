package service

import (
	"context"

	"github.com/mutition/go_start/common/metric"
	"github.com/mutition/go_start/order/adapters"
	"github.com/mutition/go_start/order/adapters/grpc"
	"github.com/mutition/go_start/order/app"
	"github.com/mutition/go_start/order/app/command"
	"github.com/mutition/go_start/order/app/query"
	"github.com/sirupsen/logrus"
	grpcclient "github.com/mutition/go_start/common/client"
)

func NewApplication(ctx context.Context) (app.Application, func() error) {
	stockClient, closeStockClient, err := grpcclient.NewGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)
	
	return newApplication(ctx,stockGRPC), func() error {
		_ = closeStockClient()
		return nil
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricclient := metric.NewTodoMetrics()
	return app.Application{
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderQueryHandler(orderRepo, logger, metricclient),
		},
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricclient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricclient),
		},
	}
}