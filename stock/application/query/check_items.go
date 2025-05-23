package query

import (
	"common/consts"
	"common/cqrs"
	"context"
	"stock/domain"
)

type CheckItemsQuery struct {
	Items []*domain.ItemWithQuantity
}

type CheckItemsResult struct {
	StockStatus consts.StockStatus
	Items       []*domain.Item
}

type CheckItemsHandler cqrs.QueryHandler[CheckItemsQuery, CheckItemsResult]

type checkItems struct {
	stockRepo domain.StockRepository
}

func (c checkItems) Handle(ctx context.Context, query CheckItemsQuery) (CheckItemsResult, error) {
	itemIDs := make([]string, len(query.Items))
	for i, item := range query.Items {
		itemIDs[i] = item.ItemID
	}

	stockItems, err := c.stockRepo.Get(ctx, itemIDs)
	if err != nil {
		return CheckItemsResult{StockStatus: consts.StockStatusUnknown}, err
	}

	for _, existing := range stockItems {
		for _, required := range query.Items {
			if existing.ItemID == required.ItemID && existing.Quantity < required.Quantity {
				return CheckItemsResult{StockStatus: consts.StockStatusInsufficient, Items: stockItems}, nil
			}
		}
	}

	return CheckItemsResult{StockStatus: consts.StockStatusSufficient, Items: stockItems}, nil
}

func NewCheckItemsHandler(stockRepo domain.StockRepository,
) CheckItemsHandler {
	return cqrs.ApplyQueryDecorator[CheckItemsQuery, CheckItemsResult](
		checkItems{stockRepo: stockRepo},
	)
}
