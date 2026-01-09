package engine

import (
	"context"
	"errors"
	"quant-trading/internal/domain/market"
	"sync"
)

type Handler interface {
	Handle(ctx context.Context, evt market.Event) error
	Types() []market.EventType
}

// Dispatcher 负责：
// 1. 事件队列
// 2. 事件 → Handler 路由
// 3. 同步 / 异步调度策略（当前：同步，保证确定性）
type Dispatcher struct {
	mu       sync.RWMutex
	handlers map[market.EventType][]Handler
	queue    chan market.Event
}

func NewDispatcher(buffer int) *Dispatcher {
	return &Dispatcher{
		handlers: make(map[market.EventType][]Handler),
		queue:    make(chan market.Event, buffer),
	}
}

// Register 注册 Handler
func (d *Dispatcher) Register(h Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, t := range h.Types() {
		d.handlers[t] = append(d.handlers[t], h)
	}
}

func (d *Dispatcher) Publish(evt market.Event) error {
	select {
	case d.queue <- evt:
		return nil
	default:
		return errors.New("事件队列已满")
	}
}

// Next 拉取下一个事件（阻塞）
func (d *Dispatcher) Next(ctx context.Context) (market.Event, error) {
	select {
	case evt := <-d.queue:
		return evt, nil
	case <-ctx.Done():
		return market.Event{}, ctx.Err()
	}
}

// Dispatch 分发事件到对应的Handler
func (d *Dispatcher) Dispatch(ctx context.Context, evt market.Event) error {
	d.mu.RLock()
	hs := d.handlers[evt.Type]
	d.mu.RUnlock()

	for _, h := range hs {
		if err := h.Handle(ctx, evt); err != nil {
			return err
		}
	}
	return nil
}
