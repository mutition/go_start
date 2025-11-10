package service

import (
	"context"

	"github.com/mutition/go_start/order/adapters"
	"github.com/mutition/go_start/order/app"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	return app.Application{
	}
}

