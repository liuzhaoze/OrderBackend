package cqrs

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

type commandDecoratorLogging[TCommand any, TResult any] struct {
	logger *logrus.Logger
	base   CommandHandler[TCommand, TResult]
}

type queryDecoratorLogging[TQuery any, TResult any] struct {
	logger *logrus.Logger
	base   QueryHandler[TQuery, TResult]
}

func (c commandDecoratorLogging[TCommand, TResult]) Handle(ctx context.Context, command TCommand) (result TResult, err error) {
	logger := c.logger.WithFields(logrus.Fields{
		"command":      getName[TCommand](command),
		"command_body": fmt.Sprintf("%+v", command),
	})

	defer func() {
		if err != nil {
			logger.Errorf("failed to execute command: %v", err)
		} else {
			logger.Info("command executed successfully")
		}
	}()

	return c.base.Handle(ctx, command)
}

func (q queryDecoratorLogging[TQuery, TResult]) Handle(ctx context.Context, query TQuery) (result TResult, err error) {
	logger := q.logger.WithFields(logrus.Fields{
		"query":      getName[TQuery](query),
		"query_body": fmt.Sprintf("%+v", query),
	})

	defer func() {
		if err != nil {
			logger.Errorf("failed to execute query: %v", err)
		} else {
			logger.Info("query executed successfully")
		}
	}()

	return q.base.Handle(ctx, query)
}
