package command

import (
	"context"

	"github.com/mutition/go_start/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
