package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"stock/application"
	"stock/application/command"
	"stock/application/query"
	"stock/infrastructure/database"
)

func NewApplication(ctx context.Context) (*application.Application, func()) {
	logger := logrus.StandardLogger()
	stockRepo := database.NewMemoryDatabase()
	return &application.Application{
			Commands: application.Commands{
				FetchItems: command.NewFetchItemsHandler(stockRepo, logger),
			},
			Queries: application.Queries{
				CheckItems: query.NewCheckItemsHandler(stockRepo, logger),
			},
		}, func() {
		}
}
