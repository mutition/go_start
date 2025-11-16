package metric

import "github.com/sirupsen/logrus"

type TodoMetrics struct{}

func (m *TodoMetrics) Inc(key string, value int) {
	logrus.Infof("Incrementing metric %s with value %d", key, value)
}

func NewTodoMetrics() *TodoMetrics {
	return &TodoMetrics{}
}
