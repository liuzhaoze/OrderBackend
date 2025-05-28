package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"stock/application"
	"stock/application/command"
	"stock/application/query"
	"stock/infrastructure/database"
)

func NewApplication(ctx context.Context) (*application.Application, func()) {
	logger := logrus.StandardLogger()
	stockRepo, err := database.NewMySQLDatabase(
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.db-name"),
	)
	if err != nil {
		logrus.Panicln(err)
	}

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
