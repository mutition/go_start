package domain

import (
	"context"

	"github.com/mutition/go_start/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error)
}
