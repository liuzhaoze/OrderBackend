package main

import (
	"common/tracing"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"order/application"
	"order/application/command"
	"order/application/query"
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
	ctx, span := tracing.StartSpan(c.Request.Context(), "Order/HTTP/POST: 创建订单")
	defer span.End()

	var requestBody ports.CreateOrderRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		resp := ports.Response{
			Data:      nil,
			ErrorCode: -1,
			Message:   err.Error(),
			TraceID:   tracing.TraceID(ctx),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if err := h.validateRequestBody(requestBody); err != nil {
		resp := ports.Response{
			Data:      nil,
			ErrorCode: -2,
			Message:   err.Error(),
			TraceID:   tracing.TraceID(ctx),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	result, err := h.app.Commands.CreateOrder.Handle(ctx, command.CreateOrderCommand{
		CustomerID: requestBody.CustomerID,
		Items:      dto.NewItemWithQuantityConverter().FromHttpBatch(requestBody.Items),
	})
	if err != nil {
		resp := ports.Response{
			Data:      nil,
			ErrorCode: -3,
			Message:   err.Error(),
			TraceID:   tracing.TraceID(ctx),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := ports.Response{
		Data: gin.H{
			"customer_id": requestBody.CustomerID,
			"order_id":    result.OrderID,
			"redirect_url": fmt.Sprintf("http://%s:%s/payment?customer-id=%s&order-id=%s",
				viper.GetString("order.http-host"), viper.GetString("order.http-port"),
				requestBody.CustomerID, result.OrderID,
			),
		},
		ErrorCode: 0,
		Message:   "success",
		TraceID:   tracing.TraceID(ctx),
	}
	c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
	ctx, span := tracing.StartSpan(c.Request.Context(), "Order/HTTP/GET: 获取订单")
	defer span.End()

	result, err := h.app.Queries.GetOrder.Handle(ctx, query.GetOrderQuery{
		OrderID:    orderId,
		CustomerID: customerId,
	})
	if err != nil {
		resp := ports.Response{
			Data:      nil,
			ErrorCode: -4,
			Message:   err.Error(),
			TraceID:   tracing.TraceID(ctx),
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
		TraceID:   tracing.TraceID(ctx),
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
