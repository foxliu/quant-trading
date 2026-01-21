package market

type Snapshot struct {
	Symbol string
	Last   float64
}

func NewSnapShot(p Price) Snapshot {
	return Snapshot{
		Symbol: p.Symbol,
		Last:   p.Last,
	}
}
