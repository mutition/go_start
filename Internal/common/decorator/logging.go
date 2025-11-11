package decorator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

type LoggingQueryDecorator[Q any, R any] struct {
	logger *logrus.Entry
	base   QueryHandler[Q, R]
}

func (d *LoggingQueryDecorator[Q, R]) Handle(ctx context.Context, cmd Q) (res R, err error) {
	logger := d.logger.WithFields(logrus.Fields{
		"query": generateQueryName(cmd),
		"query_body": fmt.Sprintf("%#v", cmd),
	})
	logger.Info("Query received")
	defer func() {
		if err != nil {
			logger.Error("Query failed ", "error", err)
		} else {
			logger.Info("Query completed")
		}
	}()
	return d.base.Handle(ctx, cmd)
}

func generateQueryName(cmd any) string {
	return reflect.TypeOf(cmd).Name()
}