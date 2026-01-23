package position

import (
	"quant-trading/internal/domain/trade"
	"time"
)

// 只读出口

type Snapshot struct {
	Symbol string
	Pos    *trade.Position
	At     time.Time
}

func NewSnapshot(pos *trade.Position) *Snapshot {
	if pos == nil {
		return nil
	}

	return &Snapshot{
		Symbol: pos.Symbol,
		At:     pos.UpdateAt,
		Pos:    pos,
	}
}

func (s *Snapshot) Name() string {
	return "position"
}

func (s *Snapshot) Timestamp() time.Time {
	return s.At
}
