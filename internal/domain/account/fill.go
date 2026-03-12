package account

type Fill struct {
	Symbol string
	Side   Side
	Qty    float64
	Price  float64
	Fee    float64
}
