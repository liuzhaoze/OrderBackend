package domain

import (
	"common/consts"
)

type Order struct {
	OrderID     string
	CustomerID  string
	Items       []*Item
	Status      consts.OrderStatus
	PaymentLink string
}
