package consumer

import (
	"context"
	"encoding/json"

	"github.com/mutition/go_start/common/broker"
	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/common/tracing"
	"github.com/mutition/go_start/payment/app"
	"github.com/mutition/go_start/payment/app/command"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Consumer struct {
	application app.Application
}

func NewConsumer(application app.Application) *Consumer {
	return &Consumer{application: application}
}

func (c *Consumer) Listen(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("failed to declare queue: %v", err)
	}
	err = ch.QueueBind(q.Name, broker.EventOrderCreated, broker.EventOrderCreated, false, nil)
	if err != nil {
		logrus.Fatalf("failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("failed to consume message from queue %s: %v", q.Name, err)
	}

	forever := make(chan struct{})
	go func() {
		for msg := range msgs {
			c.handleMessage(msg)
		}
	}()
	<-forever
	return nil
}

func (c *Consumer) handleMessage(msg amqp.Delivery) {
	logrus.Infof("received message from queue %s: %s", msg.RoutingKey, msg.Body)
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	_, span := tracing.StartSpan(ctx, "Consumer.handleMessage")
	defer span.End()
	order := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, order); err != nil {
		_ = msg.Nack(false, false)
		return
	}
	logrus.Infof("coorder: %v", order)
	if _, err := c.application.Commands.CreatePayment.Handle(ctx, command.CreatePayment{
		Order: order,
	}); err != nil {
		logrus.Infof("failed to create payment: %v", err)
		_ = msg.Nack(false, false)
		return
	}
	span.AddEvent("Payment created",trace.WithAttributes(
		attribute.String("order.id", order.Id),
		attribute.String("order.status", order.Status),
	))
	_ = msg.Ack(false)
	logrus.Infof("consumer for order %s success", order.Id)
}
