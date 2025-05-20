package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"order/ports"
)

type HttpHandler struct {
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

func (h *HttpHandler) PostCustomerCustomerIdCreate(c *gin.Context, customerId string) {
	// TODO: implement real logic
	var request ports.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(err)
	}
	fmt.Println(customerId)
	fmt.Println(request)
}

func (h *HttpHandler) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
	// TODO: implement real logic
	fmt.Println(customerId, orderId)
}
