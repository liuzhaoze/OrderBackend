package main

import (
	client "common/client/order"
	_ "common/config"
	"common/protobuf/orderpb"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	ctx      = context.Background()
	endpoint = fmt.Sprintf(
		"http://%s:%s/api",
		viper.GetString("order.http-host"),
		viper.GetString("order.http-port"),
	)
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestCreateOrderSuccess(t *testing.T) {
	customerID := "customer_test"
	requestBody := client.PostCustomerCustomerIdCreateJSONRequestBody{
		CustomerID: customerID,
		Items: []client.ItemWithQuantity{
			{
				ItemID:   "test_item_1",
				Quantity: 10,
			},
			{
				ItemID:   "test_item_2",
				Quantity: 20,
			},
			{
				ItemID:   "test_item_1",
				Quantity: 5,
			},
		},
	}

	response, err := post(t, customerID, requestBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode())
}

func post(t *testing.T, customerID string, body client.PostCustomerCustomerIdCreateJSONRequestBody) (*client.PostCustomerCustomerIdCreateResponse, error) {
	t.Helper()

	httpClient, err := client.NewClientWithResponses(endpoint)
	if err != nil {
		return nil, err
	}

	response, err := httpClient.PostCustomerCustomerIdCreateWithResponse(ctx, customerID, body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

var (
	testOrder = &orderpb.Order{
		OrderID:    "test_order_1",
		CustomerID: "test_customer_1",
		Items: []*orderpb.Item{
			{ItemID: "test_item_1", Name: "item1", Quantity: 10, PriceID: "price_test_item_1"},
			{ItemID: "test_item_2", Name: "item2", Quantity: 20, PriceID: "price_test_item_2"},
		},
		Status:      orderpb.OrderStatus_WaitingForPayment,
		PaymentLink: "test_payment_link",
	}
)

func TestUpdateOrder(t *testing.T) {
	orderGrpcClient, closeOrderGrpcClient, err := client.NewOrderGrpcClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = closeOrderGrpcClient()
	}()

	request := &orderpb.UpdateOrderRequest{
		UpdateOptions: orderpb.UpdateOption_PaymentLink,
		Order:         testOrder,
	}
	response, err := orderGrpcClient.UpdateOrder(ctx, request)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, request.UpdateOptions&orderpb.UpdateOption_PaymentLink == orderpb.UpdateOption_PaymentLink)
	assert.True(t, request.UpdateOptions&orderpb.UpdateOption_Status == orderpb.UpdateOption_Unspecified)
	assert.Equal(t, orderpb.OrderStatus_WaitingForPayment, response.Order.Status)
	assert.Equal(t, 2, len(response.Order.Items))
}
