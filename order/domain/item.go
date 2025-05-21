package domain

type Item struct {
	ItemID   string
	Name     string
	Quantity int64
	PriceID  string
}

type ItemWithQuantity struct {
	ItemID   string
	Quantity int64
}
