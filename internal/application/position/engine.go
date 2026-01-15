package position

import (
	"context"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/strategy"
	"quant-trading/pkg/utils"
)

/*
Engine
======

Position Engine 的职责：

- 接收 Planner 产生的 Order（未修正）
- 结合 Position Context
- 修正 Order.Quantity / Order.Side
*/
type Engine struct {
	ctx    *Context
	input  <-chan order.Order
	output chan<- order.Order
}

func NewEngine(ctx *Context, input <-chan order.Order, output chan<- order.Order) *Engine {
	return &Engine{
		ctx:    ctx,
		input:  input,
		output: output,
	}
}

func (e *Engine) Start(ctx context.Context) error {
	go func() {
		select {
		case <-ctx.Done():
			return
		case o := <-e.input:
			e.handle(o)
		}
	}()
	return nil
}

func (e *Engine) handle(o order.Order) {
	pos := e.ctx.Get(o.Symbol)

	// 构造一个"虚拟 signal" 用于计算
	signal := strategy.Signal{TargetQty: o.Quantity}

	delta, side := CalcDelta(pos, signal)

	if delta == 0 {
		return
	}

	o.Quantity = utils.Abs(delta)
	o.Side = side

	select {
	case e.output <- o:
	default:

	}
}
