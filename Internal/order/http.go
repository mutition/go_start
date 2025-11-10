package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mutition/go_start/order/app"
	"github.com/mutition/go_start/order/app/query"
)

type HTTPServer struct{
	app app.Application
}


func (s HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerId string) {

}

func (s HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
	order, err := s.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrderQuery{
		CustomerId: "1",
		OrderId: "1",
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order found", "order": order})
}