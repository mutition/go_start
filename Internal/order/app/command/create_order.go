package command

import (
	"context"

	"github.com/mutition/go_start/common/decorator"
	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/order/app/query"
	domain "github.com/mutition/go_start/order/domain/order"
	"github.com/sirupsen/logrus"
)

type CreateOrder struct {
	CustomerId string
	Items      []*orderpb.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderId string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
}

func NewCreateOrderHandler(orderRepo domain.Repository, stockGRPC query.StockService, 
	logger *logrus.Entry, client decorator.MetricsClient) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		&createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
		logger,
		client,
	)
}

func (h *createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	err := h.stockGRPC.CheckIfItemsInStock(ctx, cmd.Items)
	resp, err := h.stockGRPC.GetItems(ctx, []string{"123"})
	logrus.Info("response from stock service ", resp)
	var stockresponse []*orderpb.Item
	for _, item := range cmd.Items {
		stockresponse = append(stockresponse, &orderpb.Item{
			Id:       item.Id,
			Quantity: item.Quantity,
		})
	}
	o, err := h.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerId,
		Items:      stockresponse,
	})
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{
		OrderId: o.ID,
	}, nil
}
