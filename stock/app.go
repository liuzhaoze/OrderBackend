package main

import (
	"context"
	"stock/application"
	"stock/application/command"
	"stock/application/query"
	"stock/infrastructure/database"
)

func NewApplication(ctx context.Context) (*application.Application, func()) {
	stockRepo := database.NewMemoryDatabase()
	return &application.Application{
			Commands: application.Commands{
				FetchItems: command.NewFetchItemsHandler(stockRepo),
			},
			Queries: application.Queries{
				CheckItems: query.NewCheckItemsHandler(stockRepo),
			},
		}, func() {
		}
}
