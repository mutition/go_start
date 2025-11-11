package grpc

import (
	"context"

	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

type StockGRPC struct {
	client stockpb.StockServiceClient
}

func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{client: client}
}

// CheckIfItemsInStock implements query.StockService.
func (s *StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) ( error) {
	_, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{
		Items: items,
	})
	return err
}

// GetItems implements query.StockService.
func (s *StockGRPC) GetItems(ctx context.Context, itemsIds []string) ([]*orderpb.Item, error) {
	response, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{
		ItemIds: itemsIds,
	})
	logrus.Info("response from stock service ", response)
	if err != nil {
		return nil, nil
	}
	return response.Items, nil
}
