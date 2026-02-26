package runtime

import (
	"errors"
	"quant-trading/internal/application/account"
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/market"
	"quant-trading/internal/domain/strategy"
	"time"
)

/*
Runtime
=======

Runtime 表示【一个策略实例的运行时】。

设计原则：
- 一个 Runtime = 一个 Strategy
- Runtime 内部串行执行
- 并发隔离在 Dispatcher 层完成
*/
type Runtime struct {
	strategy    strategy.Strategy
	strategyCtx strategy.Context
	accountCtx  *account.Context
	bus         event.Bus
	eventCh     chan market.Event
}

func NewRuntime(s strategy.Strategy, accountCtx *account.Context, bus event.Bus, buffer int) *Runtime {
	strategyCtx := strategy.NewContext()
	strategyCtx.SetAccountContext(accountCtx)
	return &Runtime{
		strategy:    s,
		strategyCtx: strategyCtx,
		accountCtx:  accountCtx,
		bus:         bus,
		eventCh:     make(chan market.Event, buffer),
	}
}

/*
Init
----

初始化策略。
*/
func (r *Runtime) Init() error {
	if r.strategy == nil {
		return errors.New("runtime: strategy is nil")
	}
	return r.strategy.OnInit(r.strategyCtx)
}

/*
Stop
----

停止策略运行。

说明：
- 不 close eventCh（避免 Dispatcher panic）
- Runtime 停止后不再处理事件
*/
func (r *Runtime) Stop() error {
	return r.strategy.OnStop(r.strategyCtx)
}

/*
Enqueue
-------

向 Runtime 投递事件。

工程取舍：
- 非阻塞
- 队列满则丢弃
*/
func (r *Runtime) Enqueue(event market.Event) {
	select {
	case r.eventCh <- event:
	default:
		// 丢弃事件，防止慢策略拖垮系统
	}
}

/*
HandleEvent
-----------

处理单个市场事件（串行调用）。
*/
func (r *Runtime) HandleEvent(evt market.Event) ([]strategy.Signal, error) {
	// 把市场事件注入 Strategy Context 中
	r.strategyCtx.SetCurrentEvent(evt)

	signals, err := r.strategy.OnMarketEvent(r.strategyCtx, evt)
	if err != nil {
		return nil, err
	}

	// 通过 EventBus 发布 Signal （供 Risk Engine 消费）
	for _, sig := range signals {
		r.bus.Publish(&event.Envelope{
			Type:      event.EventSignal,
			Source:    "strategy-" + r.strategy.Name(),
			Timestamp: time.Now(),
			Payload:   sig,
		})
	}
	return signals, nil
}

/*
EventChan
---------

暴露只读事件通道，供 Dispatcher 消费。
*/
func (r *Runtime) EventChan() <-chan market.Event {
	return r.eventCh
}
