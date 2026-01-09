package strategyengine

import (
	"context"
	"quant-trading/internal/domain/market"
	"sync"
)

/*
Dispatcher
==========

Dispatcher 负责策略的并发调度与隔离：

- 多策略并行
- 单策略内部串行
- 行情快速入队，不阻塞上游
*/
type Dispatcher struct {
	runtimes []*Runtime

	eventCh chan market.Event
	wg      sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

func NewDispatcher(registry *Registry, buffer int) *Dispatcher {
	runtimes := make([]*Runtime, 0)
	for _, s := range registry.All() {
		runtimes = append(runtimes, NewRuntime(s))
	}

	return &Dispatcher{
		runtimes: runtimes,
		eventCh:  make(chan market.Event, buffer),
	}
}

func (d *Dispatcher) Start(parent context.Context) error {
	d.ctx, d.cancel = context.WithCancel(parent)

	for _, r := range d.runtimes {
		if err := r.Init(); err != nil {
			return err
		}
	}

	// 每个Runtime 启动一个 worker, 保证策略内串行
	for _, rt := range d.runtimes {
		rt := rt
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			d.run(rt)
		}()
	}
	return nil
}

func (d *Dispatcher) Stop() error {
	d.cancel()
	d.wg.Wait()

	for _, rt := range d.runtimes {
		_ = rt.Stop()
	}
	return nil
}

func (d *Dispatcher) Dispatch(event market.Event) {
	select {
	case d.eventCh <- event:
	default:
		// buffer 满时直接丢弃，防止拖垮系统
	}
}

func (d *Dispatcher) run(rt *Runtime) {
	for {
		select {
		case <-d.ctx.Done():
			return
		case event := <-d.eventCh:
			_, err := rt.HandleEvent(event)
			if err != nil {
				// 策略内部错误，打印错误信息，并停止策略
				_ = rt.Stop()
				return
			}
			// Signal 在下一阶段交给 Risk Engine
		}
	}
}
