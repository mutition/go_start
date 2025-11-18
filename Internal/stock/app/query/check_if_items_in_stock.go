package query

import (
	"context"

	"github.com/mutition/go_start/common/decorator"
	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/common/tracing"
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

var stub = map[string]string{
	"item-1": "price_1STfjXE7odSGW0Fa9FMkrjrq",
	"item-2": "price_1STduFE7odSGW0FaXxVaF9c3",
}

func (h *checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
	ctx, span := tracing.StartSpan(ctx, "CheckIfItemsInStock.Handle")
	defer span.End()
	var res []*orderpb.Item
	for _, item := range query.Items {
		priceId, ok := stub[item.Id]
		if !ok {
			priceId = stub["item-1"]
		}
		res = append(res, &orderpb.Item{
			Id:       item.Id,
			Quantity: item.Quantity,
			PriceId:  priceId,
		})
	}
	return res, nil
}
