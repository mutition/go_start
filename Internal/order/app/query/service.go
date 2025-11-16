package query

import (
	"context"

	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/common/genproto/stockpb"
)

type StockService interface {
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
	GetItems(ctx context.Context, itemsIds []string) ([]*orderpb.Item, error)
}
