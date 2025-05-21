package consts

type OrderStatus string

const (
	OrderStatusUnknown           OrderStatus = "UNKNOWN"
	OrderStatusPending           OrderStatus = "PENDING"
	OrderStatusWaitingForPayment OrderStatus = "WAITING_FOR_PAYMENT"
	OrderStatusPaid              OrderStatus = "PAID"
	OrderStatusFinished          OrderStatus = "FINISHED"
)
