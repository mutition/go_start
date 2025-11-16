package order

import (
	"context"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	Get(ctx context.Context, orderId string, customerId string) (*Order, error)
	Update(
		ctx context.Context, order *Order,
		updateFn func(ctx context.Context, order *Order) (*Order, error),
	) error
}

type NotFoundError struct {
	OrderId string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("order %s not found", e.OrderId)
}
