package decorator

import (
	"context"
	"fmt"
	"time"
)

type MetricsQueryDecorator[Q any, R any] struct {
	base   QueryHandler[Q, R]
	client MetricsClient
}

type MetricsClient interface {
	Inc(key string, value int)
}


func (d *MetricsQueryDecorator[Q, R]) Handle(ctx context.Context, cmd Q) (res R, err error) {
	start := time.Now()
	actionName := generateQueryName(cmd)
	defer func() {
		end := time.Since(start)
		d.client.Inc(fmt.Sprintf("query.%s.duration",actionName), int(end.Seconds()))
		if err != nil {
			d.client.Inc(fmt.Sprintf("query.%s.error",actionName), 1)
		} else {
			d.client.Inc(fmt.Sprintf("query.%s.success",actionName), 1)
		}
	}()
	return d.base.Handle(ctx, cmd)
}

