package command

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mutition/go_start/common/broker"
	"github.com/mutition/go_start/common/decorator"
	"github.com/mutition/go_start/common/tracing"
	"github.com/mutition/go_start/order/app/query"
	"github.com/mutition/go_start/order/convertor"
	domain "github.com/mutition/go_start/order/domain/order"
	"github.com/mutition/go_start/order/entity"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/sirupsen/logrus"
)

type CreateOrder struct {
	CustomerId string
	Items      []*entity.ItemWithQuantity
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
	ctx, span := tracing.StartSpan(ctx, "CreateOrder.PublishToRabbitMQ")
	defer span.End()
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


	body, err := json.Marshal(convertor.NewOrderConvertor().EntityToProto(o))
	if err != nil {
		return nil, err
	}
	header := broker.InjectRabbitMQHeaders(ctx)
	err = h.ch.PublishWithContext(ctx, broker.EventOrderCreated,broker.EventOrderCreated, false, false, amqp.Publishing{
		ContentType: "application/json",
		DeliveryMode: amqp.Persistent,
		Body:        body,
		Headers:     header,
	})
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{
		OrderId: o.ID,
	}, nil
}

func (h *createOrderHandler) validateItems(ctx context.Context, 
	items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("items are required")
	}
	items = packItems(items)
	response, err := h.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
	if err != nil {
		return nil, err
	}
	return convertor.NewItemConvertor().ProtosToEntities(response.Items), nil
}

func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}
	var result []*entity.ItemWithQuantity
	for id, quantity := range merged {
		result = append(result, &entity.ItemWithQuantity{
			ID:       id,
			Quantity: quantity,
		})
	}
	return result
}
