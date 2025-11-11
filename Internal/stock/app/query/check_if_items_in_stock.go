package query

import (
	"context"

	"github.com/mutition/go_start/common/decorator"
	"github.com/mutition/go_start/common/genproto/orderpb"
	domain "github.com/mutition/go_start/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
}

func NewCheckIfItemsInStockHandler(stockRepo domain.Repository,
	 logger *logrus.Entry, client decorator.MetricsClient) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
		&checkIfItemsInStockHandler{stockRepo: stockRepo},
		logger,
		client,
	)
}

func (h *checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var res []*orderpb.Item
	for _, item := range query.Items {
		res = append(res, &orderpb.Item{
			Id:       item.Id,
			Quantity: item.Quantity,
		})
	}
	return res, nil
}


