package command

import (
	"common/consts"
	"common/cqrs"
	"common/protobuf/orderpb"
	"common/tracing"
	"context"
	"github.com/sirupsen/logrus"
	"payment/domain"
	"payment/dto"
)

type CreatePaymentCommand struct {
	Order *domain.Order
}

type CreatePaymentResult struct {
	PaymentLink string
}

type CreatePaymentHandler cqrs.CommandHandler[CreatePaymentCommand, CreatePaymentResult]

type createPayment struct {
	paymentCreator domain.PaymentCreator
	orderGrpc      orderpb.OrderServiceClient
}

func (c createPayment) Handle(ctx context.Context, command CreatePaymentCommand) (CreatePaymentResult, error) {
	ctx, span := tracing.StartSpan(ctx, "Payment/Application/Command: create payment")
	defer span.End()

	link, err := c.paymentCreator.CreatePaymentLink(ctx, command.Order)
	if err != nil {
		return CreatePaymentResult{PaymentLink: ""}, err
	}

	orderWithPaymentLink := command.Order
	orderWithPaymentLink.Status = consts.OrderStatusWaitingForPayment
	orderWithPaymentLink.PaymentLink = link

	result, err := c.orderGrpc.UpdateOrder(ctx, &orderpb.UpdateOrderRequest{
		UpdateOptions: orderpb.UpdateOption_Status | orderpb.UpdateOption_PaymentLink,
		Order:         dto.NewOrderConverter().ToOrderGrpc(orderWithPaymentLink),
	})
	if err != nil {
		return CreatePaymentResult{PaymentLink: ""}, err
	}

	return CreatePaymentResult{PaymentLink: result.Order.PaymentLink}, nil
}

func NewCreatePaymentHandler(paymentCreator domain.PaymentCreator, orderGrpc orderpb.OrderServiceClient,
	logger *logrus.Logger,
) CreatePaymentHandler {
	return cqrs.ApplyCommandDecorator[CreatePaymentCommand, CreatePaymentResult](
		createPayment{paymentCreator: paymentCreator, orderGrpc: orderGrpc},
		logger,
	)
}
