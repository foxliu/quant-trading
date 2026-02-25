package market

import (
	"sync"
	"time"
)

// PriceView 只读视图，Strategy 永远拿不到指针
type PriceView struct {
	Symbol    string
	Last      float64
	Bid       float64
	Ask       float64
	Volume    float64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	EventTime int64
	UpdateSeq uint64
}

// MarketSnapshot 一致性快照（用于回测/决策）
type MarketSnapshot struct {
	Seq    uint64
	AsOf   int64
	Prices map[string]PriceView
}

type Context struct {
	mu     sync.RWMutex
	seq    uint64
	prices map[string]*priceState
}

type priceState struct {
	symbol    string
	last      float64
	bid       float64
	ask       float64
	volume    float64
	open      float64
	high      float64
	low       float64
	close     float64
	eventTime int64
	updateSeq uint64
}

func NewContext() *Context {
	return &Context{
		prices: make(map[string]*priceState),
	}
}

// Update 支持多 Symbol (EventBus调用)
func (c *Context) Update(symbol string, last float64, ts time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.seq++
	if ts.IsZero() {
		ts = time.Now()
	}

	ps, ok := c.prices[symbol]
	if !ok {
		ps = &priceState{symbol: symbol}
		c.prices[symbol] = ps
	}

	ps.last = last
	ps.eventTime = ts.UnixNano()
	ps.updateSeq = c.seq
}

// Get 返回只读视图
func (c *Context) Get(symbol string) (PriceView, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ps, ok := c.prices[symbol]
	if !ok {
		return PriceView{}, false
	}
	return toView(ps), true
}

// Snapshot 返回当前一致性快照（回测/决策必备）
func (c *Context) Snapshot() MarketSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snap := MarketSnapshot{
		Seq:    c.seq,
		AsOf:   time.Now().UnixNano(),
		Prices: make(map[string]PriceView, len(c.prices)),
	}
	for sym, ps := range c.prices {
		snap.Prices[sym] = toView(ps)
	}
	return snap
}

func toView(ps *priceState) PriceView {
	return PriceView{
		Symbol:    ps.symbol,
		Last:      ps.last,
		EventTime: ps.eventTime,
		UpdateSeq: ps.updateSeq,
		// Bid/Ask/OHLC 可后续扩展
	}
}
