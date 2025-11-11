package command

import (
	"context"

	"github.com/mutition/go_start/common/decorator"
	domain "github.com/mutition/go_start/order/domain/order"
	"github.com/sirupsen/logrus"
)

type UpdateOrder struct {
	Order *domain.Order
	UpdateFn func(ctx context.Context, order *domain.Order) (*domain.Order, error)
}


type UpdateOrderHandler decorator.CommandHandler[UpdateOrder,interface{}]

type updateOrderHandler struct {
	orderRepo domain.Repository
}

func NewUpdateOrderHandler(orderRepo domain.Repository, logger *logrus.Entry, client decorator.MetricsClient) UpdateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		&updateOrderHandler{orderRepo: orderRepo},
		logger,
		client,
	)
}

func (h *updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
	if cmd.UpdateFn == nil {
		logrus.Warnf("UpdateFn is nil, using default update function for order %s", cmd.Order)
		cmd.UpdateFn = func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		}
	}
	err := h.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn)
	if err != nil {
		return nil, err
	}
	return nil, nil
}