package service

import (
	"context"

	"github.com/mutition/go_start/common/broker"
	grpcclient "github.com/mutition/go_start/common/client"
	"github.com/mutition/go_start/common/metric"
	"github.com/mutition/go_start/order/adapters"
	"github.com/mutition/go_start/order/adapters/grpc"
	"github.com/mutition/go_start/order/app"
	"github.com/mutition/go_start/order/app/command"
	"github.com/mutition/go_start/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func() error) {
	stockClient, closeStockClient, err := grpcclient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)

	ch, closeCh := broker.ConnectToRabbitMQ(
		viper.GetString("rabbitmq.user"), viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"), viper.GetString("rabbitmq.port"))

	return newApplication(ctx, ch, stockGRPC), func() error {
		_ = closeCh()
		_ = ch.Close()
		_ = closeStockClient()
		return nil
	}
}

func newApplication(_ context.Context, ch *amqp.Channel, stockGRPC query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricclient := metric.NewTodoMetrics()
	return app.Application{
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderQueryHandler(orderRepo, logger, metricclient),
		},
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, ch, stockGRPC, logger, metricclient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricclient),
		},
	}
}
