package performance

import (
	"quant-trading/internal/domain/account"
	"time"
)

type EquityPoint struct {
	Time          time.Time
	Equity        float64
	Available     float64
	RealizedPnL   float64
	UnrealizedPnL float64
}

type EquityRecorder struct {
	points []EquityPoint
}

func NewEquityRecorder() *EquityRecorder {
	return &EquityRecorder{
		points: make([]EquityPoint, 0),
	}
}

func (r *EquityRecorder) Record(t time.Time, s account.Snapshot) {
	point := EquityPoint{
		Time:          t,
		Equity:        s.Balance.Frozen,
		Available:     s.Balance.Available,
		RealizedPnL:   s.RealizedPnL,
		UnrealizedPnL: s.Portfolio.Realized,
	}

	r.points = append(r.points, point)
}

func (r *EquityRecorder) Points() []EquityPoint {
	return r.points
}
