package main

import (
	client "common/client/stock"
	"common/protobuf/stockpb"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ctx = context.Background()

func TestMain(m *testing.M) {
	m.Run()
}

func TestCheckAndFetchItems(t *testing.T) {
	stockGrpcClient, closeStockGrpcClient, err := client.NewStockGrpcClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = closeStockGrpcClient() }()

	request := &stockpb.CheckAndFetchItemsRequest{Items: []*stockpb.ItemWithQuantity{
		{ItemID: "test_item_1", Quantity: 15},
		{ItemID: "test_item_2", Quantity: 20},
	}}
	response, err := stockGrpcClient.CheckAndFetchItems(ctx, request)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, stockpb.StockStatus_Insufficient, response.StatusCode)
	assert.Equal(t, 0, len(response.Items))
}
