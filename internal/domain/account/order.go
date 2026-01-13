package account

type OrderSide string

const (
	Buy  OrderSide = "BUY"
	Sell OrderSide = "SELL"
)

type Order struct {
	OrderID  string
	Symbol   string
	Side     OrderSide
	Price    float64
	Quantity int64
}
