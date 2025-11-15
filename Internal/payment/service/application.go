package service

import (
	"context"

	grpcclient "github.com/mutition/go_start/common/client"
	"github.com/mutition/go_start/common/metric"
	"github.com/mutition/go_start/payment/adapters"
	"github.com/mutition/go_start/payment/app"
	"github.com/mutition/go_start/payment/app/command"
	"github.com/mutition/go_start/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func() error) {
	orderClient, closeOrderClient, err := grpcclient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	memoryProcessor := processor.NewInmenProcessor()
	return newApplication(ctx, orderGRPC, memoryProcessor), func() error {
		return closeOrderClient()
	}
}

func newApplication(ctx context.Context, orderGRPC command.OrderService,
	processor *processor.InmenProcessor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricclient := metric.NewTodoMetrics()
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricclient),
		},
	}
}
