package main

import (
	client "common/client/order"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"payment/application"
	"payment/application/command"
	"payment/infrastructure/stripe"
)

func NewApplication(ctx context.Context) (*application.Application, func()) {
	logger := logrus.StandardLogger()
	stripeKey := viper.GetString("STRIPE_KEY")
	if stripeKey == "" {
		logger.Panicln("empty stripe key, please set STRIPE_KEY environment variable")
	}
	paymentCreator := stripe.NewStripe(stripeKey, "TODO:")
	orderGrpcClient, closeOrderGrpcClient, err := client.NewOrderGrpcClient(ctx)
	if err != nil {
		logger.Panicln(err)
	}

	return &application.Application{
			Commands: application.Commands{
				CreatePayment: command.NewCreatePaymentHandler(paymentCreator, orderGrpcClient, logger),
			},
			Queries: application.Queries{},
		}, func() {
			_ = closeOrderGrpcClient()
		}
}
