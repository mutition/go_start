package query

import (
	"context"

	"github.com/mutition/go_start/common/decorator"
	"github.com/mutition/go_start/common/genproto/orderpb"
	domain "github.com/mutition/go_start/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type GetItems struct {
	ItemIds []string
}

type GetItemsHandler decorator.QueryHandler[GetItems, []*orderpb.Item]

type getItemsHandler struct {
	stockRepo domain.Repository
}

func NewGetItemsHandler(stockRepo domain.Repository,
	logger *logrus.Entry, client decorator.MetricsClient) GetItemsHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorators[GetItems, []*orderpb.Item](
		&getItemsHandler{stockRepo: stockRepo},
		logger,
		client,
	)
}

func (h *getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*orderpb.Item, error) {
	items, err := h.stockRepo.GetItems(ctx, query.ItemIds)
	if err != nil {
		return nil, err
	}
	return items, nil
}
