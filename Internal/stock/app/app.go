package app

import "github.com/mutition/go_start/stock/app/query"

type Application struct {
	Commands Commads
	Queries  Queries
}

type Commads struct {
}

type Queries struct {
	GetItems            query.GetItemsHandler
	CheckIfItemsInStock query.CheckIfItemsInStockHandler
}
