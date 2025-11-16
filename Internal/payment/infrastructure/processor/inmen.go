package processor

import (
	"context"

	"github.com/mutition/go_start/common/genproto/orderpb"
)

type InmenProcessor struct {
}

func NewInmenProcessor() *InmenProcessor {
	return &InmenProcessor{}
}

func (p *InmenProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	return "inmen_payment_link", nil
}
