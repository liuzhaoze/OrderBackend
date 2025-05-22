package consts

type StockStatus int32

const (
	StockStatusUnknown StockStatus = iota
	StockStatusSufficient
	StockStatusInsufficient
)
