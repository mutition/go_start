package app

import "github.com/mutition/go_start/order/app/query"

type Application struct {
	Commands Commads
	Queries  Queries
}

type Commads struct {
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderQueryHandler
}