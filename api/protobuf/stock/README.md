# stock 模块的 gRPC 接口定义

## CheckAndGetItemsFromStock

> order 模块调用该接口获得物品的 Name 和 PriceID
>
> order 模块的 http POST 请求只有物品的 ItemID 和数量

该接口接受物品和数量信息，检查库存中是否有足够的物品

通过 `StatusCode` 指示库存是否充足

如果库存充足，则扣减库存，并返回剩余库存

如果库存不足，则返回当前库存数量
