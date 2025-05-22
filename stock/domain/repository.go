package domain

import (
	"context"
	"fmt"
)

type StockRepository interface {
	Get(ctx context.Context, itemIDs []string) ([]*Item, error)
	Update(ctx context.Context, items []*Item, updateFunc func(context.Context, []*Item) ([]*Item, error)) ([]*Item, error)
}

type NotFoundError struct {
	MissingItemIDs []string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Item not found with IDs: %v", e.MissingItemIDs)
}
