package stock

import (
	"context"
	"fmt"
	"strings"

	"github.com/mutition/go_start/common/genproto/orderpb"
)

type Repository interface {
	GetItems(ctx context.Context, itemIds []string) ([]*orderpb.Item, error)
}

type NotFoundError struct {
	MissingItemIds []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("these items %s not found in stock", strings.Join(e.MissingItemIds, ", "))
}