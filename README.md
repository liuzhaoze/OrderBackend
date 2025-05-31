# Order Backend

## Overview

![architecture](./docs/architecture.drawio.svg)

1. order 模块接受 HTTP POST 创建订单请求，返回一个指向支付页面的重定向连接；该页面向 order 模块发送 HTTP GET 请求获取订单状态
2. order 模块调用 stock 模块的 gRPC 接口检查和扣减库存
3. 扣减库存成功后，order 模块创建一个 `PENDING` 状态的订单，并向消息队列发送 `order.created` 事件
4. payment 模块消费 `order.created` 事件，使用 Stripe SDK 创建支付意图
5. payment 模块调用 order 模块的 gRPC 接口更新订单状态 `PENDING -> WAITING_FOR_PAYMENT`，并且更新支付链接
6. 待用户从支付页面点击支付链接并完成支付后，Stripe CLI 会向 payment 模块的 HTTP Webhook 发送支付成功事件
7. payment 接收到支付成功事件后，向消息队列广播 `order.paid` 事件
8. order 模块消费 `order.paid` 事件，更新订单状态 `WAITING_FOR_PAYMENT -> PAID`
9. process 模块消费 `order.paid` 事件，处理订单（例如发货）
10. process 模块处理完成订单后调用 order 模块的 gRPC 接口更新订单状态 `PAID -> FINISHED`

## Deployment

1. 安装 [air](https://github.com/air-verse/air) 用于热重载 `go install github.com/air-verse/air@latest`
2. 启动容器 `docker compose up --build -d`
3. 启动 Stripe CLI 并监听支付成功事件 `stripe listen --forward-to localhost:8301/webhook`
4. 启动 stock 模块
   - 打开新终端
   - 执行 `cd stock && air .`
5. 启动 order 模块
   - 打开新终端
   - 执行 `cd order && air .`
6. 启动 payment 模块
   - 打开新终端
   - 设置环境变量 `export STRIPE_KEY=<StripeWeb沙盒中的API密钥>`
   - 设置环境变量 `export STRIPE_ENDPOINT_SECRET=<StripeCLI中的Webhook签名密钥>`
   - 执行 `cd payment && air .`
7. 启动 process 模块
   - 打开新终端
   - 执行 `cd process && air .`

## Testing

向 `http://localhost:8101/api/customer/{customerID}/create` 发送 HTTP POST 请求创建订单

请求体：

```json
{
    "customerID": "{{customerID}}",
    "items": [
        {
            "itemID": "prod_SNjRcpjxpiazxk",
            "quantity": 5
        },
        {
            "itemID": "prod_SNjQpQjNC8QuaD",
            "quantity": 5
        },
        {
            "itemID": "prod_SNjRcpjxpiazxk",
            "quantity": 15
        }
    ]
}
```

得到相应：

```json
{
    "data": {
        "customer_id": "0a69556f-fe95-48eb-bb0a-099d9b0f2a50",
        "order_id": "68ce52f4cdbf4a483de63446",
        "redirect_url": "http://127.0.0.1:8101/payment?customer-id=0a69556f-fe95-48eb-bb0a-099d9b0f2a50&order-id=68ce52f4cdbf4a483de63446"
    },
    "errorCode": 0,
    "message": "success",
    "traceID": "654f54b90d7ed003c4d05fad87e5a925"
}
```

打开 `redirect_url`，点击支付链接完成支付

测试卡号 `4242 4242 4242 4242`，其他信息任意

打开 zipkin Web UI `http://localhost:8600` 查看链路追踪

右上角输入 `traceID` 查询

## Domain Driven Design

![DDD](./docs/DDD.drawio.svg)

---

使用 `go work` 管理多模块开发。
