package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"order/application"
	"order/application/command"
	"order/application/query"
	"order/domain"
	"order/dto"
	"order/ports"
)

type HttpHandler struct {
	app *application.Application
}

func NewHttpHandler(app *application.Application) *HttpHandler {
	return &HttpHandler{app: app}
}

func (h *HttpHandler) PostCustomerCustomerIdCreate(c *gin.Context, customerId string) {
	var requestBody ports.CreateOrderRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		resp := ports.Response{
			Data:      nil,
			ErrorCode: -1,
			Message:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if err := h.validateRequestBody(requestBody); err != nil {
		resp := ports.Response{
			Data:      nil,
			ErrorCode: -2,
			Message:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	convertedItems := make([]*domain.ItemWithQuantity, len(requestBody.Items))
	for i, item := range requestBody.Items {
		convertedItems[i] = dto.NewItemWithQuantityConverter().FromHttp(item)
	}
	result, err := h.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrderCommand{
		CustomerID: requestBody.CustomerID,
		Items:      convertedItems,
	})
	if err != nil {
		resp := ports.Response{
			Data:      nil,
			ErrorCode: -3,
			Message:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := ports.Response{
		Data: gin.H{
			"customer_id": requestBody.CustomerID,
			"order_id":    result.OrderID,
		},
		ErrorCode: 0,
		Message:   "success",
	}
	c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
	result, err := h.app.Queries.GetOrder.Handle(c.Request.Context(), query.GetOrderQuery{
		OrderID:    orderId,
		CustomerID: customerId,
	})
	if err != nil {
		resp := ports.Response{
			Data:      nil,
			ErrorCode: -4,
			Message:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := ports.Response{
		Data: gin.H{
			"customer_id":  result.Order.CustomerID,
			"order_id":     result.Order.OrderID,
			"status":       result.Order.Status,
			"payment_link": result.Order.PaymentLink,
		},
		ErrorCode: 0,
		Message:   "success",
	}
	c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) validateRequestBody(requestBody ports.CreateOrderRequest) error {
	for _, item := range requestBody.Items {
		if item.Quantity <= 0 {
			return errors.New(fmt.Sprintf("item quantity must be a positive integer, got %d from %s", item.Quantity, item.ItemID))
		}
	}
	return nil
}
