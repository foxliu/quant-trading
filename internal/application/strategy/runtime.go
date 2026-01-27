package strategyengine

import (
	"errors"
	"quant-trading/internal/application/account"
	dAccount "quant-trading/internal/domain/account"
	"quant-trading/internal/domain/market"
	"quant-trading/internal/domain/strategy"
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
	Account dAccount.Descriptor

	strategy    strategy.Strategy
	strategyCtx strategy.Context
	accountCtx  *account.Context

	eventCh chan market.Event
}

func NewRuntime(s strategy.Strategy, accountCtx *account.Context, buf int) *Runtime {
	ctx := strategy.NewContext()
	ctx.SetAccountContext(accountCtx)
	return &Runtime{
		strategy:    s,
		strategyCtx: ctx,
		accountCtx:  accountCtx,
		eventCh:     make(chan market.Event, buf),
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
func (r *Runtime) HandleEvent(event market.Event) ([]strategy.Signal, error) {
	return r.strategy.OnMarketEvent(r.strategyCtx, event)
}

/*
EventChan
---------

暴露只读事件通道，供 Dispatcher 消费。
*/
func (r *Runtime) EventChan() <-chan market.Event {
	return r.eventCh
}
