package pnl

import (
	"sync"
	"time"
)

type Context struct {
	mu sync.Mutex

	symbol string

	// === Position Mirror ===
	qty      int64
	avgPrice float64

	// === PnL ===
	realized   float64
	unrealized float64

	markPrice float64
	updateAt  time.Time
}

func NewContext(symbol string) *Context {
	return &Context{symbol: symbol}
}
