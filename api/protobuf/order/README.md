# order 模块的 gRPC 接口定义

## UpdateOrder

- payment 模块调用该接口将 order 的状态从 PENDING 更新为 WAITING_FOR_PAYMENT 并且更新 payment link
- process 模块调用该接口将 order 的状态从 PAID 更新为 READY
