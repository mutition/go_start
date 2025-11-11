package query

import (
	"context"

	"github.com/mutition/go_start/common/genproto/orderpb"
)

type StockService interface {
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) ( error)
	GetItems(ctx context.Context, itemsIds []string) ([]*orderpb.Item, error)
}

