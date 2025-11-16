package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func ConnectToRabbitMQ(user, password, host, port string) (*amqp.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	conn, err := amqp.Dial(address)
	if err != nil {
		logrus.Errorf("failed to dial to rabbitmq: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		logrus.Errorf("failed to dial to rabbitmq: %v", err)
	}
	err = ch.ExchangeDeclare(EventOrderCreated, "direct", true, false, false, false, nil)
	if err != nil {
		logrus.Errorf("failed to dial to rabbitmq: %v", err)
	}
	err = ch.ExchangeDeclare(EventOrderPaid, "direct", true, false, false, false, nil)
	if err != nil {
		logrus.Errorf("failed to dial to rabbitmq: %v", err)
	}
	return ch, conn.Close
}
