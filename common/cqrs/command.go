package cqrs

import "context"

type CommandHandler[TCommand any, TResult any] interface {
	Handle(ctx context.Context, command TCommand) (TResult, error)
}

func ApplyCommandDecorator[TCommand any, TResult any](
	handler CommandHandler[TCommand, TResult],
) CommandHandler[TCommand, TResult] {
	// TODO: logger, metrics
	return handler
}
