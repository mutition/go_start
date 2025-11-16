package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mutition/go_start/common/genproto/orderpb"
	"github.com/mutition/go_start/order/app"
	"github.com/mutition/go_start/order/app/command"
	"github.com/mutition/go_start/order/app/query"
)

type HTTPServer struct {
	app app.Application
}

func (s HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerId string) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := s.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
		CustomerId: req.CustomerId,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order created",
		"order_id": order.OrderId, "customer_id": customerId,
		"redirect": fmt.Sprintf("http://localhost:8282/payment/success?customerID=%s&orderID=%s", customerId, order.OrderId)})
}

func (s HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
	order, err := s.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		CustomerId: customerId,
		OrderId:    orderId,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Order found",
		"data": map[string]any{
			"id":           order.ID,
			"customer_id":  order.CustomerID,
			"status":       order.Status,
			"payment_link": order.PaymentLink,
			"items":        order.Items,
		},
	})
}
