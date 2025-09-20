package database

import (
	"common/tracing"
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/rand"
	"stock/domain"
	"time"
)

type MySQLDatabase struct {
	db *gorm.DB
}

func NewMySQLDatabase(user, password, host, port, dbName string) (*MySQLDatabase, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &MySQLDatabase{db: db}, nil
}

type itemModel struct {
	ItemID    string
	Name      string
	Quantity  int64
	PriceID   string
	Version   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (i itemModel) TableName() string {
	return "items"
}

func (m *MySQLDatabase) Get(ctx context.Context, itemIDs []string) ([]*domain.Item, error) {
	ctx, span := tracing.StartSpan(ctx, "Stock Repository: get")
	defer span.End()

	var result []*itemModel
	m.db.Where("item_id IN ?", itemIDs).Find(&result)

	return m.unmarshalItemBatch(result), nil
}

func (m *MySQLDatabase) Update(ctx context.Context, items []*domain.Item, updateFunc func(context.Context, []*domain.Item) ([]*domain.Item, error)) ([]*domain.Item, error) {
	ctx, span := tracing.StartSpan(ctx, "Stock Repository: update")
	defer span.End()

	// return m.updatePessimistic(ctx, items, updateFunc)
	return m.updateOptimistic(ctx, items, updateFunc)
}

func (m *MySQLDatabase) updatePessimistic(ctx context.Context, items []*domain.Item, updateFunc func(context.Context, []*domain.Item) ([]*domain.Item, error)) ([]*domain.Item, error) {
	var (
		updatedItems []*domain.Item
		err          error
	)

	err = m.db.Transaction(func(tx *gorm.DB) (err error) {
		itemIDs := make([]string, len(items))
		for i, item := range items {
			itemIDs[i] = item.ItemID
		}

		var result []*itemModel
		if err = tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("item_id IN ?", itemIDs).Find(&result).Error; err != nil {
			return
		}

		if updatedItems, err = updateFunc(ctx, m.unmarshalItemBatch(result)); err != nil {
			return
		}

		for _, item := range m.marshalItemBatch(updatedItems) {
			if err = tx.Where("item_id = ?", item.ItemID).Select("quantity").Updates(item).Error; err != nil {
				// 如果不添加 .Select() ，gorm 会只更新非零值字段，无法将某些字段更新为零值
				// 参见 https://gorm.io/docs/update.html#Updates-Changed-Only
				return
			}
		}

		return nil
	})

	return updatedItems, err
}

// 使用抖动指数退避避免惊群效应
const (
	maxRetries     = 5
	baseDelayMs    = 5
	maxDelayMs     = 100
	jitterFactorMs = 20
)

var ErrOptimisticLockConflict = errors.New("optimistic lock conflict")

func (m *MySQLDatabase) updateOptimistic(ctx context.Context, items []*domain.Item, updateFunc func(context.Context, []*domain.Item) ([]*domain.Item, error)) ([]*domain.Item, error) {
	var (
		updatedItems []*domain.Item
		err          error
	)

	for attempt := 0; attempt < maxRetries; attempt++ {
		updatedItems, err = m.tryOptimisticUpdate(ctx, items, updateFunc)

		if err == nil {
			return updatedItems, nil
		}

		if !errors.Is(err, ErrOptimisticLockConflict) {
			return nil, err
		}

		// 发生冲突时，指数退避延迟后重试
		if attempt < maxRetries-1 {
			waitWithJitterBackoff(attempt)
		}
	}

	return nil, fmt.Errorf("optimistic lock failed after %d attempts: %w", maxRetries, err)
}

// waitWithJitterBackoff 带抖动的退避算法，避免惊群效应
func waitWithJitterBackoff(attempt int) {
	// 指数退避
	exponentialDelay := min(baseDelayMs*(1<<attempt), maxDelayMs)

	// 添加随机抖动，范围在 0 到 jitterFactorMs 之间
	jitter := rand.Intn(jitterFactorMs)

	totalDelay := time.Duration(exponentialDelay+jitter) * time.Millisecond
	time.Sleep(totalDelay)
}

func (m *MySQLDatabase) tryOptimisticUpdate(ctx context.Context, items []*domain.Item, updateFunc func(context.Context, []*domain.Item) ([]*domain.Item, error)) ([]*domain.Item, error) {
	var (
		updatedItems []*domain.Item
		err          error
	)

	err = m.db.Transaction(func(tx *gorm.DB) (err error) {
		itemIDs := make([]string, len(items))
		for i, item := range items {
			itemIDs[i] = item.ItemID
		}

		// 读取当前数据和版本号
		var result []*itemModel
		if err = tx.Where("item_id IN ?", itemIDs).Find(&result).Error; err != nil {
			return
		}

		// 保存原始版本号
		originalVersions := make(map[string]int64)
		for _, item := range result {
			originalVersions[item.ItemID] = item.Version
		}

		if updatedItems, err = updateFunc(ctx, m.unmarshalItemBatch(result)); err != nil {
			return
		}

		for _, item := range m.marshalItemBatch(updatedItems) {
			originalVersion := originalVersions[item.ItemID]

			// 使用版本号作为条件进行更新，同时递增版本号
			info := tx.Model(&itemModel{}).
				Where("item_id = ? AND version = ?", item.ItemID, originalVersion).
				Updates(map[string]any{
					"quantity": item.Quantity,
					"version":  originalVersion + 1,
				})

			if info.Error != nil {
				return info.Error
			}

			if info.RowsAffected == 0 {
				// 如果没有满足条件的行被更新，说明版本号不匹配，在本次事务中要更改的数据已被其他事务修改
				// 返回冲突错误用于重试
				return ErrOptimisticLockConflict
			}
		}

		return nil
	})

	return updatedItems, err
}

func (m *MySQLDatabase) marshalItem(item *domain.Item) *itemModel {
	return &itemModel{
		ItemID:   item.ItemID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (m *MySQLDatabase) marshalItemBatch(items []*domain.Item) []*itemModel {
	result := make([]*itemModel, len(items))
	for i, item := range items {
		result[i] = m.marshalItem(item)
	}
	return result
}

func (m *MySQLDatabase) unmarshalItem(item *itemModel) *domain.Item {
	return &domain.Item{
		ItemID:   item.ItemID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (m *MySQLDatabase) unmarshalItemBatch(items []*itemModel) []*domain.Item {
	result := make([]*domain.Item, len(items))
	for i, item := range items {
		result[i] = m.unmarshalItem(item)
	}
	return result
}
