package broker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
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

type RabbitMQHeaderCarrier map[string]interface{}

func (h RabbitMQHeaderCarrier) Get(key string) string {
	value, ok := h[key]
	if !ok {
		return ""
	}
	return value.(string)
}

func (h RabbitMQHeaderCarrier) Set(key, value string) {
	h[key] = value
}

func (h RabbitMQHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(h))
	for key := range h {
		keys = append(keys, key)
	}
	return keys
}

func InjectRabbitMQHeaders(ctx context.Context)map[string]interface{} {
	carrier := make(RabbitMQHeaderCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier
}

func ExtractRabbitMQHeaders(ctx context.Context,headers map[string]interface{})context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, RabbitMQHeaderCarrier(headers))
}