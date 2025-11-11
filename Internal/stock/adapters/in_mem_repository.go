package adapters

import (
	"context"
	"sync"

	"github.com/mutition/go_start/common/genproto/orderpb"
	domain "github.com/mutition/go_start/stock/domain/stock"
)

type MemoryStockRepository struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

func NewMemoryStockRepository() *MemoryStockRepository {
	return &MemoryStockRepository{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}

var stub = map[string]*orderpb.Item{
	"item-1": {
		Id:       "item-1",
		Name:     "Item-1",
		PriceId:  "price1",
		Quantity: 10,
	},
	"item-2": {
		Id:       "item-2",
		Name:     "Item-2",
		PriceId:  "price2",
		Quantity: 20,
	},
	"item-3": {
		Id:       "item-3",
		Name:     "Item-3",
		PriceId:  "price3",
		Quantity: 30,
	},
}

func (m *MemoryStockRepository) GetItems(ctx context.Context, itemIds []string) ([]*orderpb.Item, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var res []*orderpb.Item
	var missingItemIds []string
	for _, itemId := range itemIds {
		item, ok := m.store[itemId]
		if ok {
			res = append(res, item)
		} else {
			missingItemIds = append(missingItemIds, itemId)
		}
	}
	if len(res) == len(itemIds) {
		return res, nil
	}
	return res, domain.NotFoundError{MissingItemIds: missingItemIds}
}
