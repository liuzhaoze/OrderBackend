#

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

## Domain Driven Design

![DDD](./docs/DDD.drawio.svg)

---

使用 `go work` 管理多模块开发。
