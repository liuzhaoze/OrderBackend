# order 模块的 OpenAPI 定义

## `/customer/{customer_id}/orders/{order_id}`

### `GET /customer/{customer_id}/orders/{order_id}`

获取用户 ID为 `customer_id` 的订单 ID 为 `order_id` 的订单信息。

## `/customer/{customer_id}/orders`

### `POST /customer/{customer_id}/orders`

为用户 ID 为 `customer_id` 的用户创建一个新的订单。

request body 包括用户 ID 和用户购买的商品和数量。
