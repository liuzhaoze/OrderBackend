package database

import (
	"context"
	"crypto/md5"
	"fmt"
	"order/domain"
	"strconv"
	"sync"
	"time"
)

type MemoryDatabase struct {
	lock *sync.RWMutex
	db   []*domain.Order
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{lock: &sync.RWMutex{}, db: make([]*domain.Order, 0)}
}

func (m *MemoryDatabase) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	orderID := fmt.Sprintf("%x", md5.Sum([]byte(order.CustomerID+strconv.Itoa(int(time.Now().Unix())))))
	newOrder, err := domain.NewOrder(orderID, order.CustomerID, order.Items, order.Status, order.PaymentLink)
	if err != nil {
		return nil, err
	}

	m.db = append(m.db, newOrder)
	return newOrder, nil
}

func (m *MemoryDatabase) Get(ctx context.Context, orderID string, customerID string) (*domain.Order, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, o := range m.db {
		if o.OrderID == orderID && o.CustomerID == customerID {
			return domain.NewOrder(o.OrderID, o.CustomerID, o.Items, o.Status, o.PaymentLink)
		}
	}

	return nil, &domain.NotFoundError{
		OrderID:    orderID,
		CustomerID: customerID,
	}
}

// Update 首先找到 order 在数据库中对应的对象，然后按照 updateFunc 的逻辑修改该对象，最后返回修改后的对象
func (m *MemoryDatabase) Update(ctx context.Context, order *domain.Order, updateFunc func(context.Context, *domain.Order) (*domain.Order, error)) (*domain.Order, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, o := range m.db {
		if o.OrderID == order.OrderID && o.CustomerID == order.CustomerID {
			if updatedOrder, err := updateFunc(ctx, o); err != nil {
				return nil, err
			} else {
				return domain.NewOrder(updatedOrder.OrderID, updatedOrder.CustomerID, updatedOrder.Items, updatedOrder.Status, updatedOrder.PaymentLink)
			}
		}
	}

	return nil, &domain.NotFoundError{
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
	}
}
