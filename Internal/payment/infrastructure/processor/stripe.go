package processor

import (
	"context"
	"encoding/json"

	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/checkout/session"
)

type StripeProcessor struct {
	apiKey string
}

func NewStripeProcessor(apiKey string) *StripeProcessor {
	if apiKey == "" {
		panic("apiKey is required")
	}
	stripe.Key = apiKey
	return &StripeProcessor{apiKey: apiKey}
}

func (p *StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	var items []*stripe.CheckoutSessionLineItemParams
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			//先写死
			Price:    stripe.String("price_1STduFE7odSGW0FaXxVaF9c3"),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return "", err
	}
	metadata := map[string]string{
		"order_id":    order.Id,
		"customer_id": order.CustomerId,
		"items":       string(itemsJSON),
		"status":      order.Status,
	}
	params := &stripe.CheckoutSessionParams{
		Metadata:   metadata,
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

var successURL = "http://localhost:8282/payment/success"
