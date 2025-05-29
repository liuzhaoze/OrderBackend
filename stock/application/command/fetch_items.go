package command

import (
	"common/cqrs"
	"common/tracing"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sort"
	"stock/domain"
	"strings"
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

func (f fetchItems) Handle(ctx context.Context, command FetchItemsCommand) (FetchItemsResult, error) {
	ctx, span := tracing.StartSpan(ctx, "Stock/Application/Command: fetch items")
	defer span.End()

	// 使用 Redis 分布式锁确保并发安全
	lockKey := getRedisLockKey(command)
	lockValue := uuid.New().String()
	if ok, err := f.locker.Lock(ctx, lockKey, lockValue); err != nil || !ok {
		return FetchItemsResult{Items: nil}, fmt.Errorf("failed to acquire lock %s: %v", lockKey, err)
	}
	defer func() {
		if ok, err := f.locker.Unlock(ctx, lockKey, lockValue); err != nil || !ok {
			logrus.Errorf("failed to release lock %s: %v", lockKey, err)
		}
	}()

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

func NewFetchItemsHandler(stockRepo domain.StockRepository, locker domain.Locker,
	logger *logrus.Logger,
) FetchItemsHandler {
	return cqrs.ApplyCommandDecorator[FetchItemsCommand, FetchItemsResult](
		fetchItems{stockRepo: stockRepo, locker: locker},
		logger,
	)
}

const (
	redisLockPrefix = "Stock:Command:FetchItems:"
)

func getRedisLockKey(cmd FetchItemsCommand) string {
	itemIDs := make([]string, len(cmd.Items))
	for i, item := range cmd.Items {
		itemIDs[i] = item.ItemID
	}
	sort.Strings(itemIDs)
	return redisLockPrefix + strings.Join(itemIDs, ".")
}
