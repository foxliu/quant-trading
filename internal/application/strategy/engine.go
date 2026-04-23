package strategyengine

import (
	"context"
	"quant-trading/internal/domain/market"
	"sync"
)

/*
Engine
======

Strategy Engine 的对外门面（Facade）。

它解决三个问题：
1. 系统何时开始 / 停止跑策略
2. 行情从哪里进入策略系统
3. Engine 与外部模块的边界隔离

Engine 不直接执行策略逻辑，
所有实际计算均委托给 Dispatcher / Runtime。
*/
type Engine struct {
	// ctx 用于整体生命周期控制（优雅退出、系统停止）
	ctx        context.Context
	cancel     context.CancelFunc
	dispatcher *Dispatcher

	mu      *sync.Mutex
	started bool
}

func NewEngine(dispatcher *Dispatcher) *Engine {
	ctx, cancel := context.WithCancel(context.Background())
	return &Engine{
		ctx:        ctx,
		cancel:     cancel,
		dispatcher: dispatcher,
		mu:         &sync.Mutex{},
	}
}

func (e *Engine) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.started {
		return nil
	}

	if err := e.dispatcher.Start(e.ctx); err != nil {
		return err
	}

	e.started = true
	return nil
}

func (e *Engine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.started {
		return nil
	}

	e.cancel()
	if err := e.dispatcher.Stop(); err != nil {
		return err
	}
	e.started = false
	return nil
}

// OnMarketEvent 是 Strategy Engine 的唯一行情入口
func (e *Engine) OnMarketEvent(event market.Event) {
	e.mu.Lock()
	started := e.started
	e.mu.Unlock()

	if !started {
		return
	}
	e.dispatcher.Dispatch(event)
}
