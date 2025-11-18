package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	common "github.com/mutition/go_start/common"
	client "github.com/mutition/go_start/common/client/order"
	"github.com/mutition/go_start/common/tracing"
	"github.com/mutition/go_start/order/app"
	"github.com/mutition/go_start/order/app/command"
	"github.com/mutition/go_start/order/app/query"
	"github.com/mutition/go_start/order/convertor"
)

type HTTPServer struct {
	app app.Application
	common.BaseResponse
}

func (s HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerId string) {
	ctx, span := tracing.StartSpan(c.Request.Context(), "PostCustomerCustomerIdOrders")
	defer span.End()
	var (
		err error
		req  client.CreateOrderRequest
		resp struct {
			CustomerId string `json:"customer_id"`
			OrderId string `json:"order_id"`
			Redirect string `json:"redirect"`
		}
	)
	defer s.Response(c, err, &resp)

	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	order, err := s.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerId: customerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}
	//traceID := tracing.TraceID(ctx)
	resp.CustomerId = customerId
	resp.OrderId = order.OrderId
	resp.Redirect = fmt.Sprintf("http://localhost:8282/payment/success?customerID=%s&orderID=%s", customerId, order.OrderId)
	}

func (s HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
	var (
		err error
		resp struct {
			Order *client.Order `json:"order"`
		}
	)
	defer s.Response(c, err, &resp)
	ctx, span := tracing.StartSpan(c.Request.Context(), "GetCustomerCustomerIdOrdersOrderId")
	defer span.End()
	order, err := s.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerId: customerId,
		OrderId:    orderId,
	})
	if err != nil {
		return
	}
	resp.Order = convertor.NewOrderConvertor().EntityToClient(order)
}
