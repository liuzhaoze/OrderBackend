package main

import (
	"common/metrics"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"stock/application"
	"stock/application/command"
	"stock/application/query"
	"stock/infrastructure/database"
	"stock/infrastructure/lock"
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

	locker := lock.NewRedisLocker(
		viper.GetString("redis.host"),
		viper.GetString("redis.port"),
		viper.GetDuration("redis.lock.expiration"),
		viper.GetInt("redis.lock.retry-number"),
		viper.GetDuration("redis.lock.retry-delay"),
	)

	metricsClient := metrics.GetPrometheusClient()

	return &application.Application{
			Commands: application.Commands{
				FetchItems: command.NewFetchItemsHandler(stockRepo, locker, logger, metricsClient),
			},
			Queries: application.Queries{
				CheckItems: query.NewCheckItemsHandler(stockRepo, logger, metricsClient),
			},
		}, func() {
		}
}
