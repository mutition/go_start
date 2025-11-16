package service

import (
	"context"

	grpcclient "github.com/mutition/go_start/common/client"
	"github.com/mutition/go_start/common/metric"
	"github.com/mutition/go_start/payment/adapters"
	"github.com/mutition/go_start/payment/app"
	"github.com/mutition/go_start/payment/app/command"
	"github.com/mutition/go_start/payment/domain"
	"github.com/mutition/go_start/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func() error) {
	orderClient, closeOrderClient, err := grpcclient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, stripeProcessor), func() error {
		return closeOrderClient()
	}
}

func newApplication(_ context.Context, orderGRPC command.OrderService,
	processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricclient := metric.NewTodoMetrics()
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricclient),
		},
	}
}
