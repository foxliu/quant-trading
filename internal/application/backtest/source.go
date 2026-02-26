package backtest

import (
	"context"
	"quant-trading/internal/application/event"
)

type BacktestSource struct {
	events []*event.Envelope
	index  int
}

func NewBacktestSource(events []*event.Envelope) *BacktestSource {
	return &BacktestSource{events: events}
}

func (s *BacktestSource) Next(ctx context.Context) (*event.Envelope, bool) {
	if s.index >= len(s.events) {
		return nil, false
	}
	evt := s.events[s.index]
	s.index++
	return evt, true
}
