package main

import (
	client "common/client/order"
	_ "common/config"
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
