package application

import (
	"order/application/command"
	"order/application/query"
)

type Commands struct {
	CreateOrder command.CreateOrderHandler
	UpdateOrder command.UpdateOrderHandler
}

type Queries struct {
	GetOrder query.GetOrderHandler
}

type Application struct {
	Commands Commands
	Queries  Queries
}
