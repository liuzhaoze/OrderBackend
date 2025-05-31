package cqrs

import (
	"common/metrics"
	"context"
	"github.com/sirupsen/logrus"
)

type CommandHandler[TCommand any, TResult any] interface {
	Handle(ctx context.Context, command TCommand) (TResult, error)
}

func ApplyCommandDecorator[TCommand any, TResult any](
	handler CommandHandler[TCommand, TResult],
	logger *logrus.Logger,
	metricsClient metrics.Client,
) CommandHandler[TCommand, TResult] {
	return commandDecoratorLogging[TCommand, TResult]{
		logger: logger,
		base: commandDecoratorMetrics[TCommand, TResult]{
			client: metricsClient,
			base:   handler,
		},
	}
}
