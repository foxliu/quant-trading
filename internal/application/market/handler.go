package market

import "time"

/*
Handler 是外部行情系统的唯一入口
*/
type Handler struct {
	ctx *Context
}

func NewHandler(ctx *Context) *Handler {
	return &Handler{ctx: ctx}
}

func (h *Handler) OnTick(symbol string, price float64, ts time.Time) {
	h.ctx.Update(symbol, price, ts)
}
