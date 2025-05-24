package cqrs

import (
	"context"
	"github.com/sirupsen/logrus"
)

type QueryHandler[TQuery any, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}

func ApplyQueryDecorator[TQuery any, TResult any](
	handler QueryHandler[TQuery, TResult],
	logger *logrus.Logger,
) QueryHandler[TQuery, TResult] {
	// TODO: metrics
	return queryDecoratorLogging[TQuery, TResult]{
		logger: logger,
		base:   handler,
	}
}
