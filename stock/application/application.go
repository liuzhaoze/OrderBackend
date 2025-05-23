package application

import (
	"stock/application/command"
	"stock/application/query"
)

type Commands struct {
	FetchItems command.FetchItemsHandler
}

type Queries struct {
	CheckItems query.CheckItemsHandler
}

type Application struct {
	Commands Commands
	Queries  Queries
}
