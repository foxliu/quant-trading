package backtest

import (
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/order"
	"time"
)

type Simulator struct {
	priceProvider func(symbol string) float64
}

func NewSimulator(priceProvider func(symbol string) float64) *Simulator {
	return &Simulator{priceProvider: priceProvider}
}

func (s *Simulator) Execute(order order.Order) ([]execution.Fill, error) {
	price := s.priceProvider(order.Symbol())

	fill := execution.Fill{
		OrderID:   order.ID(),
		AccountID: order.AccountID(),
		Symbol:    order.Symbol(),
		Side:      order.Side(),
		Qty:       order.Qty(),
		Price:     price,
		Fee:       0,
		Timestamp: time.Now(),
	}
	return []execution.Fill{fill}, nil
}
