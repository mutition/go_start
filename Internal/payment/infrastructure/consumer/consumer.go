package consumer

import (
	"github.com/mutition/go_start/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Listen(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("failed to consume message from queue %s: %v", q.Name, err)
	}

	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.handleMessage(msg, q, ch)
		}
	}()
	<-forever
	return nil
}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
	logrus.Infof("received message from queue %s: %s", msg.RoutingKey, msg.Body)
	_ = msg.Ack(false)
}