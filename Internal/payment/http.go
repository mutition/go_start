package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mutition/go_start/common/broker"
	"github.com/mutition/go_start/common/genproto/orderpb"
	domain "github.com/mutition/go_start/payment/domain"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type PaymentHandler struct {
	ch *amqp.Channel
}

func NewPaymentHandler(ch *amqp.Channel) PaymentHandler {
	return PaymentHandler{ch}
}

func (h PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", h.handleWebhook)
}

func (h PaymentHandler) handleWebhook(c *gin.Context) {
	logrus.Info("Webhook received from stripe")

	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error("Error reading request body: ", err)
		c.JSON(http.StatusServiceUnavailable, err)
		return
	}
	event, err := webhook.ConstructEventWithOptions(body, c.Request.Header.Get("Stripe-Signature"),
		viper.GetString("stripe-webhook-secret"), webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true})

	if err != nil {
		logrus.Error("Error constructing event: ", err)
		c.JSON(http.StatusServiceUnavailable, err)
		return
	}

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			logrus.Error("Error unmarshalling event data: ", err)
			c.JSON(http.StatusServiceUnavailable, err)
			return
		}
		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			logrus.Info("Checkout session completed and paid")
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			var items []*orderpb.Item
			err = json.Unmarshal([]byte(session.Metadata["items"]), &items)
			if err != nil {
				logrus.Error("Error unmarshalling items: ", err)
				c.JSON(http.StatusServiceUnavailable, err)
				return
			}
			marshalledOrder, err := json.Marshal(&domain.Order{
				ID:          session.Metadata["order_id"],
				CustomerID:  session.Metadata["customer_id"],
				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
				Items:       items,
				PaymentLink: session.Metadata["payment_link"],
			})
			if err != nil {
				logrus.Error("Error marshalling order: ", err)
				c.JSON(http.StatusServiceUnavailable, err)
				return
			}

			_ = h.ch.PublishWithContext(ctx, broker.EventOrderPaid, broker.EventOrderPaid, false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        marshalledOrder,
			})
			logrus.Infof("Order %s paid to %s, body: %s", session.Metadata["order_id"], session.Metadata["customer_id"], marshalledOrder)
		}
	}
	c.JSON(http.StatusOK, nil)
}
