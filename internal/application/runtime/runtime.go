package runtime

import (
	"context"
	"errors"
	"quant-trading/internal/application/account"
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/market"
	"quant-trading/internal/domain/strategy"
	"quant-trading/internal/infrastructure/logger"
	"sync"
	"time"

	"go.uber.org/zap"
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

	wg     sync.WaitGroup
	done   chan struct{}
	logger *zap.Logger
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
		done:        make(chan struct{}),
		logger:      logger.Logger.With(zap.String("module", "runtime.runtime")),
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
Start
-----

启动策略运行时（必须调用）。

流程：
1. 调用策略 OnInit
2. 启动后台协程消费事件（串行执行 OnMarketEvent）
3. 支持 context 取消
*/
func (r *Runtime) Start(ctx context.Context) error {
	if r.strategy == nil {
		return errors.New("runtime: strategy is nil")
	}

	// 初始化策略
	if err := r.strategy.OnInit(r.strategyCtx); err != nil {
		return err
	}
	r.wg.Add(1)
	go r.runEventLook(ctx)

	r.logger.Info("Runtime 已启动", zap.String("strategy", r.strategy.Name()))
	return nil
}

// runEventLoop 后台事件消费协程（串行处理）
func (r *Runtime) runEventLook(ctx context.Context) {
	defer r.wg.Done()

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("Strategy 外部调用停止", zap.String("strategy", r.strategy.Name()))
			return
		case <-r.done:
			r.logger.Info("Strategy 已停止", zap.String("strategy", r.strategy.Name()))
			return
		case evt, ok := <-r.eventCh:
			if !ok {
				return
			}
			// 串行处理（避免并发问题）
			if _, err := r.HandleEvent(evt); err != nil {
				r.logger.Error("Strategy 处理事件失败", zap.Error(err), zap.String("strategy", r.strategy.Name()))
			}
		}
	}
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
	close(r.done)
	r.wg.Wait()

	if err := r.strategy.OnStop(r.strategyCtx); err != nil {
		r.logger.Error("Strategy 停止失败", zap.Error(err), zap.String("strategy", r.strategy.Name()))
		return err
	}

	r.logger.Info("Strategy Runtime 已停止", zap.String("strategy", r.strategy.Name()))
	return nil
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

	// 通过 EventBus 发布 Signal （供 RiskContext Engine 消费）
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

func (r *Runtime) GetStrategyContext() strategy.Context {
	return r.strategyCtx
}
