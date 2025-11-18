package command

import (
	"context"

	"github.com/mutition/go_start/common/decorator"
	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/common/tracing"
	"github.com/mutition/go_start/payment/domain"
	"github.com/sirupsen/logrus"
)

type CreatePayment struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]

type createPaymentHandler struct {
	//processor for stripe
	processor domain.Processor
	orderGRPC OrderService
}

func NewCreatePaymentHandler(processor domain.Processor, orderGRPC OrderService,
	logger *logrus.Entry, metricclient decorator.MetricsClient) CreatePaymentHandler {
	return decorator.ApplyCommandDecorators[CreatePayment, string](
		&createPaymentHandler{processor: processor, orderGRPC: orderGRPC},
		logger,
		metricclient,
	)
}

func (h *createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "CreatePayment.CreatePayment")
	defer span.End()
	link, err := h.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}
	logrus.Infof("payment for order %s created || create_payment || link: %s", cmd.Order.Id, link)
	newOrder := &orderpb.Order{
		Id:          cmd.Order.Id,
		CustomerId:  cmd.Order.CustomerId,
		Status:      "waiting_for_payment",
		Items:       cmd.Order.Items,
		PaymentLink: link,
	}
	err = h.orderGRPC.UpdateOrder(ctx, newOrder)
	return link, err
}
