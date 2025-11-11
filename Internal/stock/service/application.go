package service

import (
	"context"

	"github.com/mutition/go_start/stock/adapters"
	"github.com/mutition/go_start/stock/app"
	"github.com/mutition/go_start/stock/app/query"
	"github.com/sirupsen/logrus"
	"github.com/mutition/go_start/common/metric"
)

func NewApplication(ctx context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	logger := logrus.NewEntry(logrus.New())
	client := metric.NewTodoMetrics()
	return app.Application{
		Queries: app.Queries{
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, client),
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, client),
		},
	}
}
