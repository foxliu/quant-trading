package portfolio

type PositionSnapshot struct {
	Symbol    string
	Quantity  float64
	AvgPrice  float64
	LastPrice float64
}

type Snapshot struct {
	Positions []PositionSnapshot
	Realized  float64
}
