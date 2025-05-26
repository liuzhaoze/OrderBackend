package application

import "payment/application/command"

type Commands struct {
	CreatePayment command.CreatePaymentHandler
}

type Queries struct {
}

type Application struct {
	Commands Commands
	Queries  Queries
}
