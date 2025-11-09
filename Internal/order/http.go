package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mutition/go_start/order/app"
)

type HTTPServer struct{
	app app.Application
}


func (s HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerId string) {

}

func (s HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
}