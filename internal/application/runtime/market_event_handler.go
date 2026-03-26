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

	// 关键修正：正确断言为 market.Event（与 HandleEvent 签名完全一致）
	if marketEvt, ok := evt.Payload.(market.Event); ok {
		_, err := h.rt.HandleEvent(marketEvt)
		if err != nil {
			return
		} // ← 直接调用，无需任何改动
	}
}
