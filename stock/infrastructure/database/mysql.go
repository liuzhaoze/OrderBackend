package database

import (
	"common/tracing"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	return m.updatePessimistic(ctx, items, updateFunc)
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
