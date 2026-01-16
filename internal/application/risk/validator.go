package risk

import "quant-trading/internal/domain/order"

func (e *engine) handle(o order.Order) {
	if !e.validate(o) {
		return
	}

	o = e.clamp(o)

	select {
	case e.output <- o:
	default:

	}
}

func (e *engine) validate(o order.Order) bool {
	if o.Quantity <= 0 {
		return false
	}

	if o.Side == "" {
		return false
	}
	return true
}

func (e *engine) clamp(o order.Order) order.Order {
	// 单笔最大下单量限制
	if e.ctx.MaxOrderQty > 0 && o.Quantity > e.ctx.MaxPositionQty {
		o.Quantity = e.ctx.MaxOrderQty
	}
	// 仓位限制在 v1 不做精确校验（需要当前仓位）
	// 这里只做结构预留
	return o
}
