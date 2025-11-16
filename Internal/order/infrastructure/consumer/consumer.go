package consumer

import (
	"context"
	"encoding/json"

	"github.com/mutition/go_start/common/broker"
	"github.com/mutition/go_start/order/app"
	"github.com/mutition/go_start/order/app/command"
	domain "github.com/mutition/go_start/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	application app.Application
}

func NewConsumer(application app.Application) *Consumer {
	return &Consumer{application: application}
}

func (c *Consumer) Listen(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(broker.EventOrderPaid, true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("failed to declare queue: %v", err)
	}
	err = ch.QueueBind(q.Name, broker.EventOrderPaid, broker.EventOrderPaid, false, nil)
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
	order := &domain.Order{}
	if err := json.Unmarshal(msg.Body, order); err != nil {
		_ = msg.Nack(false, false)
		return
	}

	if _, err := c.application.Commands.UpdateOrder.Handle(context.TODO(), command.UpdateOrder{
		Order: order,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.IsPaid(); err != nil {
				return nil, err
			}
			return order, nil
		},
	}); err != nil {
		logrus.Infof("order %s failed to update payment: %v", order.ID, err)
		return
	}
	_ = msg.Ack(false)
	logrus.Infof("consumer for order %s success", order.ID)
}
