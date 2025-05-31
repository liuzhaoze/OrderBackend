package cqrs

import (
	"common/metrics"
	"context"
	"github.com/sirupsen/logrus"
)

type QueryHandler[TQuery any, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}

func ApplyQueryDecorator[TQuery any, TResult any](
	handler QueryHandler[TQuery, TResult],
	logger *logrus.Logger,
	metricsClient metrics.Client,
) QueryHandler[TQuery, TResult] {
	return queryDecoratorLogging[TQuery, TResult]{
		logger: logger,
		base: queryDecoratorMetrics[TQuery, TResult]{
			client: metricsClient,
			base:   handler,
		},
	}
}
