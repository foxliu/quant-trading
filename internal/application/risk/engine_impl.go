package risk

import (
	"context"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/strategy"
	"sync"
)

type engine struct {
	ctx    *Context
	input  <-chan order.Order
	output chan<- order.Order

	stop chan struct{}
}

func NewEngine(ctx *Context, input <-chan order.Order, output chan<- order.Order) Engine {
	return &engine{
		ctx:    ctx,
		input:  input,
		output: output,
		stop:   make(chan struct{}),
	}
}

func (e *engine) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-e.stop:
				return
			case o := <-e.input:
				e.handle(o)
			}
		}
	}()
}

func (e *engine) Stop() error {
	close(e.stop)
	return nil
}

func (e *engine) Consume(o order.Order) {
	select {
	case e.input <- o:
	default:
		// 丢弃， 保证不阻塞上游
	}
}
