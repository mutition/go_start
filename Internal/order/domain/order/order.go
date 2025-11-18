package order

import (
	"errors"

	"github.com/mutition/go_start/order/entity"
)

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*entity.Item
}

func NewOrder(id, customerID, status, paymentLink string, items []*entity.Item) (*Order, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if customerID == "" {
		return nil, errors.New("customerID is required")
	}
	if status == "" {
		return nil, errors.New("status is required")
	}
	if items == nil {
		return nil, errors.New("items are required")
	}
	return &Order{
		ID:          id,
		CustomerID:  customerID,
		Status:      status,
		Items:       items,
		PaymentLink: paymentLink,
	}, nil
}


func (o *Order) IsPaid() error {
	if o.Status != "paid" {
		return errors.New("order is not paid, order id: " + o.ID)
	}
	return nil
}
