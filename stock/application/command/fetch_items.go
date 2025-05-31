package command

import (
	"common/cqrs"
	"common/metrics"
	"common/tracing"
	"context"
	"fmt"
	"github.com/google/uuid"
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
	locker    domain.Locker
}

const redisLockPrefix = "stock_lock:"

func (f fetchItems) Handle(ctx context.Context, command FetchItemsCommand) (FetchItemsResult, error) {
	ctx, span := tracing.StartSpan(ctx, "Stock/Application/Command: fetch items")
	defer span.End()

	query := make([]*domain.Item, len(command.Items))
	for i, item := range command.Items {
		query[i] = &domain.Item{ItemID: item.ItemID} // 只需要赋值 ItemID 用于指定要更新的 item 即可
	}

	// 对每个 item 加锁
	lockValue := uuid.New().String()
	for _, item := range query {
		lockKey := redisLockPrefix + item.ItemID
		if ok, err := f.locker.Lock(ctx, lockKey, lockValue); err != nil || !ok {
			return FetchItemsResult{Items: nil}, fmt.Errorf("failed to acquire lock for item %s: %v", item.ItemID, err)
		}
	}
	defer func() {
		for _, item := range query {
			lockKey := redisLockPrefix + item.ItemID
			if ok, err := f.locker.Unlock(ctx, lockKey, lockValue); err != nil || !ok {
				logrus.Errorf("failed to release lock for item %s: %v", item.ItemID, err)
			}
		}
	}()

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

func NewFetchItemsHandler(stockRepo domain.StockRepository, locker domain.Locker,
	logger *logrus.Logger,
	metricsClient metrics.Client,
) FetchItemsHandler {
	return cqrs.ApplyCommandDecorator[FetchItemsCommand, FetchItemsResult](
		fetchItems{stockRepo: stockRepo, locker: locker},
		logger,
		metricsClient,
	)
}
