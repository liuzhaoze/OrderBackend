package cqrs

import "context"

type QueryHandler[TQuery any, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}

func ApplyQueryDecorator[TQuery any, TResult any](
	handler QueryHandler[TQuery, TResult],
) QueryHandler[TQuery, TResult] {
	// TODO: logger, metrics
	return handler
}
