package capital

type Snapshot struct {
	Total     float64
	Available float64
	Frozen    float64

	FrozenOrders map[string]float64
}
