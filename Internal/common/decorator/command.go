package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CommandHandler[C any, R any] interface {
	Handle(ctx context.Context, query C) (R, error)
}

func ApplyCommandDecorators[C any, R any](handler CommandHandler[C, R], logger *logrus.Entry, client MetricsClient) CommandHandler[C, R] {
	return &LoggingQueryDecorator[C, R]{
		logger: logger,
		base: &MetricsQueryDecorator[C, R]{
			base:   handler,
			client: client,
		},
	}
}
