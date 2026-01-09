// Package domain portfolio 资产组合
package portfolio

// Portfolio 账户模型
type Portfolio struct {
	Cash   float64 // 可用资金
	Equity float64 // 总权益（含浮盈、浮亏）
	Margin float64 // 已占用保证金

	Positions map[string]*Position // key: Instrument.ID
}

func NewPortfolio(initialCash float64) *Portfolio {
	return &Portfolio{
		Cash:      initialCash,
		Equity:    initialCash,
		Margin:    0,
		Positions: make(map[string]*Position),
	}
}
