package cqrs

import (
	"common/metrics"
	"context"
	"time"
)

type commandDecoratorMetrics[TCommand any, TResult any] struct {
	client metrics.Client
	base   CommandHandler[TCommand, TResult]
}

type queryDecoratorMetrics[TQuery any, TResult any] struct {
	client metrics.Client
	base   QueryHandler[TQuery, TResult]
}

func (c commandDecoratorMetrics[TCommand, TResult]) Handle(ctx context.Context, command TCommand) (result TResult, err error) {
	start := time.Now()
	name := getName[TCommand](command)

	defer func() {
		duration := time.Since(start)
		c.client.RecordTime(name, float64(duration.Milliseconds()))
		if err != nil {
			c.client.CountCall(name, "fail")
		} else {
			c.client.CountCall(name, "success")
		}
	}()

	return c.base.Handle(ctx, command)
}

func (q queryDecoratorMetrics[TQuery, TResult]) Handle(ctx context.Context, query TQuery) (result TResult, err error) {
	start := time.Now()
	name := getName[TQuery](query)

	defer func() {
		duration := time.Since(start)
		q.client.RecordTime(name, float64(duration.Milliseconds()))
		if err != nil {
			q.client.CountCall(name, "fail")
		} else {
			q.client.CountCall(name, "success")
		}
	}()

	return q.base.Handle(ctx, query)
}
