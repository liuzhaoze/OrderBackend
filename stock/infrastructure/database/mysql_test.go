// 执行测试之前需要先执行 reset_test_db.sql 重置测试数据库
package database

import (
	_ "common/config"
	"common/tracing"
	"context"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"stock/domain"
	"sync"
	"testing"
)

func TestMySQLDatabase_Update(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	shutdown, err := tracing.OTelTracer(
		viper.GetString("zipkin.host"),
		viper.GetString("zipkin.port"),
		"TestMySQLDatabase_Update",
	)
	assert.NoError(t, err)
	defer func() {
		_ = shutdown(ctx)
	}()

	stockRepo, err := NewMySQLDatabase(
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.db-name")+"_test",
	)
	assert.NoError(t, err)

	// 存储查询结果
	var stockItems []*domain.Item

	// 要扣减库存的商品
	demand := []*domain.ItemWithQuantity{
		{ItemID: "prod_SNjRcpjxpiazxk", Quantity: 10},
	}

	stockItems, err = stockRepo.Get(ctx, []string{demand[0].ItemID})
	assert.NoError(t, err)
	assert.NotEmpty(t, stockItems, "stockItems should not be empty")
	initialQuantity := stockItems[0].Quantity

	var (
		wg                    sync.WaitGroup
		nConcurrentGoRoutines int64 = 10
	)
	for range nConcurrentGoRoutines {
		wg.Add(1)
		go func() {
			defer wg.Done()

			query := make([]*domain.Item, len(demand))
			for i, item := range demand {
				query[i] = &domain.Item{ItemID: item.ItemID}
			}

			_, err = stockRepo.Update(ctx, query, func(c context.Context, items []*domain.Item) ([]*domain.Item, error) {
				for _, item := range items {
					for _, required := range demand {
						if item.ItemID != required.ItemID {
							continue
						}
						item.Quantity -= required.Quantity
					}
				}
				return items, nil
			})
			assert.NoError(t, err)
		}()
	}

	wg.Wait()

	stockItems, err = stockRepo.Get(ctx, []string{demand[0].ItemID})
	assert.NoError(t, err)
	assert.NotEmpty(t, stockItems, "stockItems should not be empty")
	finalQuantity := stockItems[0].Quantity

	expectedQuantity := initialQuantity - nConcurrentGoRoutines*demand[0].Quantity
	assert.Equal(t, expectedQuantity, finalQuantity)
}
