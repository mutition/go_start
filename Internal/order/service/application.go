package service

import (
	"context"

	"github.com/mutition/go_start/common/metric"
	"github.com/mutition/go_start/order/adapters"
	"github.com/mutition/go_start/order/app"
	"github.com/mutition/go_start/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		FullTimestamp:    false,
		ForceColors:      true,
		DisableColors:    false,
	})
	logger := logrus.NewEntry(logrusLogger)
	metricclient := metric.NewTodoMetrics()
	return app.Application{
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderQueryHandler(orderRepo, logger, metricclient),
		},
	}
}
