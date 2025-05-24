package cqrs

import (
	"context"
	"github.com/sirupsen/logrus"
)

type CommandHandler[TCommand any, TResult any] interface {
	Handle(ctx context.Context, command TCommand) (TResult, error)
}

func ApplyCommandDecorator[TCommand any, TResult any](
	handler CommandHandler[TCommand, TResult],
	logger *logrus.Logger,
) CommandHandler[TCommand, TResult] {
	// TODO: metrics
	return commandDecoratorLogging[TCommand, TResult]{
		logger: logger,
		base:   handler,
	}
}
