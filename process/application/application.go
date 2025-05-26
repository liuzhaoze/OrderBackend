package application

import "process/application/command"

type Commands struct {
	ProcessOrder command.ProcessOrderHandler
}

type Queries struct {
}

type Application struct {
	Commands Commands
	Queries  Queries
}
