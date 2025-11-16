package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

func ApplyQueryDecorators[Q any, R any](handler QueryHandler[Q, R], logger *logrus.Entry, client MetricsClient) QueryHandler[Q, R] {
	return &LoggingQueryDecorator[Q, R]{
		logger: logger,
		base: &MetricsQueryDecorator[Q, R]{
			base:   handler,
			client: client,
		},
	}
}
