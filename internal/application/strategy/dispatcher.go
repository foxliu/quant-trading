package strategyengine

import (
	"context"
	"quant-trading/internal/application/risk"
	"quant-trading/internal/domain/market"
	"sync"
)

/*
Dispatcher
==========

Dispatcher 负责【事件广播 + Runtime 并发调度】。

工程语义（已冻结）：

- Dispatcher 不创建 Runtime，只调度 Runtime

- 每个 Runtime 内部串行

- 多 Runtime 并行

- 每个市场事件会被广播给所有 Runtime
*/
type Dispatcher struct {
	runtimes []*Runtime
	risk     risk.Engine

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup
}

/*
NewDispatcher
-------------

Dispatcher 是一个“调度器”，而不是“运行时工厂”。

- Runtime 必须在外部创建（Account-aware）
- Dispatcher 只接收 Runtime 集合
*/
func NewDispatcher(runtimes []*Runtime, riskEngine risk.Engine) *Dispatcher {
	return &Dispatcher{
		runtimes: runtimes,
		risk:     riskEngine,
	}
}

/*
Start
-----

启动 Dispatcher：

- 初始化 Dispatcher 生命周期
- 初始化所有 Runtime
- 为每个 Runtime 启动一个独立 worker
*/
func (d *Dispatcher) Start(parent context.Context) error {
	d.ctx, d.cancel = context.WithCancel(parent)

	if d.risk != nil {
		if err := d.risk.Start(d.ctx); err != nil {
			return err
		}
	}

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

/*
Stop
----

停止 Dispatcher：

- 先 cancel Dispatcher context
- 等待所有 Runtime worker 退出
- 再调用 Runtime.Stop 做资源清理
*/
func (d *Dispatcher) Stop() error {
	if d.cancel != nil {
		d.cancel()
	}

	d.wg.Wait()

	for _, rt := range d.runtimes {
		_ = rt.Stop()
	}

	if d.risk != nil {
		_ = d.risk.Stop()
	}
	return nil
}

/*
Dispatch
--------

将市场事件广播给所有 Runtime。

工程取舍：
- 不阻塞上游
- Runtime 自身负责处理慢 / 丢事件的策略
*/
func (d *Dispatcher) Dispatch(event market.Event) {
	for _, rt := range d.runtimes {
		rt.Enqueue(event)
	}
}

/*
run
---

单个 Runtime 的执行循环：

- 串行消费事件
- Runtime 出错只影响自身
*/
func (d *Dispatcher) run(rt *Runtime) {
	for {
		select {
		case <-d.ctx.Done():
			return
		case event := <-rt.EventChan():
			signals, err := rt.HandleEvent(event)
			if err != nil {
				// 策略内部错误，打印错误信息，并停止策略
				_ = rt.Stop()
				return
			}
			// Signal 在下一阶段交给 Risk Engine
			// TODO: 实现signal消费逻辑
			_ = signals
			// for _, sig := range signals {
			//     d.risk.Consume(sig)
			// }
		}
	}
}
