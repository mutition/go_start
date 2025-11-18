package domain

import (
	"context"

	"github.com/mutition/go_start/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error)
}

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
}


func ToProto(o *Order) *orderpb.Order {
	return &orderpb.Order{
		Id:          o.ID,
		CustomerId:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       o.Items,
	}
}
