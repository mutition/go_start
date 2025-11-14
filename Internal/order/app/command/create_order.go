package command

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mutition/go_start/common/broker"
	"github.com/mutition/go_start/common/decorator"
	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/order/app/query"
	domain "github.com/mutition/go_start/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"

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
	ch *amqp.Channel
}

func NewCreateOrderHandler(orderRepo domain.Repository,ch *amqp.Channel, 
	stockGRPC query.StockService,logger *logrus.Entry, 
	client decorator.MetricsClient) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	if ch == nil {
		panic("ch is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		&createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC, ch: ch},
		logger,
		client,
	)
}

func (h *createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	responseItems, err := h.validateItems(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}
	o, err := h.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerId,
		Items:      responseItems,
	})
	if err != nil {
		return nil, err
	}

	q, err := h.ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}


	body, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	err = h.ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		DeliveryMode: amqp.Persistent,
		Body:        body,
	})
	if err != nil {
		return nil, err
	}

	

	return &CreateOrderResult{
		OrderId: o.ID,
	}, nil
}

func (h *createOrderHandler) validateItems(ctx context.Context, 
	items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("items are required")
	}
	items = packItems(items)
	response, err := h.stockGRPC.CheckIfItemsInStock(ctx, items)
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}

func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.Id] += item.Quantity
	}
	var result []*orderpb.ItemWithQuantity
	for id, quantity := range merged {
		result = append(result, &orderpb.ItemWithQuantity{
			Id:       id,
			Quantity: quantity,
		})
	}
	return result
}
