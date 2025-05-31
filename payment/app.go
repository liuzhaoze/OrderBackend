package main

import (
	client "common/client/order"
	"common/metrics"
	"context"
	"fmt"
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
	paymentBaseURL := fmt.Sprintf("http://%s:%s", viper.GetString("order.http-host"), viper.GetString("order.http-port"))
	paymentCreator := stripe.NewStripe(stripeKey, paymentBaseURL)
	orderGrpcClient, closeOrderGrpcClient, err := client.NewOrderGrpcClient(ctx)
	if err != nil {
		logger.Panicln(err)
	}

	metricsClient := metrics.GetPrometheusClient()

	return &application.Application{
			Commands: application.Commands{
				CreatePayment: command.NewCreatePaymentHandler(paymentCreator, orderGrpcClient, logger, metricsClient),
			},
			Queries: application.Queries{},
		}, func() {
			_ = closeOrderGrpcClient()
		}
}
