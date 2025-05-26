package database

import (
	"context"
	"stock/domain"
	"sync"
)

type MemoryDatabase struct {
	lock *sync.RWMutex
	db   []*domain.Item
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		lock: &sync.RWMutex{},
		db: []*domain.Item{
			{ItemID: "prod_SNjRcpjxpiazxk", Name: "Pencil", Quantity: 100, PriceID: "price_1RSxyuPSQHt2xYB8XhSJRSVX"},
			{ItemID: "prod_SNjQpQjNC8QuaD", Name: "Book", Quantity: 200, PriceID: "price_1RSxy2PSQHt2xYB8uZzd0XQx"},
		},
	}
}

func (m *MemoryDatabase) Get(ctx context.Context, itemIDs []string) ([]*domain.Item, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var (
		result  = make([]*domain.Item, 0, len(itemIDs))
		missing = make([]string, 0)
	)

	for _, target := range itemIDs {
		isExist := false
		for _, existing := range m.db {
			if existing.ItemID == target {
				isExist = true
				result = append(result, domain.NewItem(existing.ItemID, existing.Name, existing.Quantity, existing.PriceID))
			}
		}
		if !isExist {
			missing = append(missing, target)
		}
	}

	if len(missing) > 0 {
		return nil, &domain.NotFoundError{MissingItemIDs: missing}
	}
	return result, nil
}

// Update 首先找到 items 在数据库中对应的所有对象，然后将这些对象按照 updateFunc 的逻辑进行修改，最后返回修改后的对象们
func (m *MemoryDatabase) Update(ctx context.Context, items []*domain.Item, updateFunc func(context.Context, []*domain.Item) ([]*domain.Item, error)) ([]*domain.Item, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	var (
		result  = make([]*domain.Item, 0, len(items))
		missing = make([]string, 0)
	)

	for _, target := range items {
		isExist := false
		for _, existing := range m.db {
			if existing.ItemID == target.ItemID {
				isExist = true
				result = append(result, existing)
			}
		}
		if !isExist {
			missing = append(missing, target.ItemID)
		}
	}

	if len(missing) > 0 {
		return nil, &domain.NotFoundError{MissingItemIDs: missing}
	}

	if updatedItems, err := updateFunc(ctx, result); err != nil {
		return nil, err
	} else {
		newUpdatedItems := make([]*domain.Item, 0, len(updatedItems))
		for _, item := range updatedItems {
			newUpdatedItems = append(newUpdatedItems, domain.NewItem(item.ItemID, item.Name, item.Quantity, item.PriceID))
		}
		return newUpdatedItems, nil
	}
}
