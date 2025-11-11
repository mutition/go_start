package query

import (
	"context"

	"github.com/mutition/go_start/common/decorator"
	domain "github.com/mutition/go_start/order/domain/order"
	"github.com/sirupsen/logrus"
)

type GetCustomerOrder struct {
	CustomerId string
	OrderId    string
}

type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

type getCustomerOrderQueryHandler struct {
	orderRepo domain.Repository
}


func NewGetCustomerOrderQueryHandler(orderRepo domain.Repository,
	logger *logrus.Entry, client decorator.MetricsClient) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
		&getCustomerOrderQueryHandler{orderRepo: orderRepo},
		logger,
		client,
	)
}

func (h *getCustomerOrderQueryHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
	o, err := h.orderRepo.Get(ctx, query.OrderId, query.CustomerId)
	if err != nil {
		return nil, err
	}
	return o, nil
}
