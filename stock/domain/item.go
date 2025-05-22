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

func NewItem(itemID string, name string, quantity int64, priceID string) *Item {
	return &Item{ItemID: itemID, Name: name, Quantity: quantity, PriceID: priceID}
}
