package runtime

import (
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/market"
)

// marketEventHandler 为每个 Runtime 注册的专用 Handler
type marketEventHandler struct {
	rt *Runtime
}

func (h *marketEventHandler) Handle(evt *event.Envelope) {
	if evt.Type != event.EventMarketPrice {
		return
	}

	// 只入队，由 Runtime 后台协程串行消费，避免并发直接执行策略。
	if marketEvt, ok := evt.Payload.(market.Event); ok {
		h.rt.Enqueue(marketEvt)
	}
}
