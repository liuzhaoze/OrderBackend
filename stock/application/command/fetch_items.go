package command

import (
	"common/cqrs"
	"context"
	"github.com/sirupsen/logrus"
	"stock/domain"
)

type FetchItemsCommand struct {
	Items []*domain.ItemWithQuantity
}

type FetchItemsResult struct {
	Items []*domain.Item
}

type FetchItemsHandler cqrs.CommandHandler[FetchItemsCommand, FetchItemsResult]

type fetchItems struct {
	stockRepo domain.StockRepository
}

func (f fetchItems) Handle(ctx context.Context, command FetchItemsCommand) (FetchItemsResult, error) {
	query := make([]*domain.Item, len(command.Items))
	for i, item := range command.Items {
		query[i] = &domain.Item{ItemID: item.ItemID} // 只需要赋值 ItemID 用于指定要更新的 item 即可
	}

	remaining, err := f.stockRepo.Update(ctx, query, func(c context.Context, items []*domain.Item) ([]*domain.Item, error) {
		for _, item := range items {
			for _, required := range command.Items {
				if item.ItemID != required.ItemID {
					continue
				}
				item.Quantity -= required.Quantity
			}
		}
		return items, nil
	})
	if err != nil {
		return FetchItemsResult{Items: nil}, err
	}
	return FetchItemsResult{Items: remaining}, nil
}

func NewFetchItemsHandler(stockRepo domain.StockRepository,
	logger *logrus.Logger,
) FetchItemsHandler {
	return cqrs.ApplyCommandDecorator[FetchItemsCommand, FetchItemsResult](
		fetchItems{stockRepo: stockRepo},
		logger,
	)
}
