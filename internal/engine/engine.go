package engine

import (
	"context"
	"errors"
)

// Engine 只负责：
// 1. 生命周期
// 2. 主循环
// 3. 错误隔离
type Engine struct {
	dispatcher *Dispatcher
	started    bool
}

func NewEngine(dispatcher *Dispatcher) *Engine {
	return &Engine{
		dispatcher: dispatcher,
	}
}

func (e *Engine) Run(ctx context.Context) error {
	if e.started {
		return errors.New("engine 已经在运行中")
	}
	e.started = true

	for {
		evt, err := e.dispatcher.Next(ctx)
		if err != nil {
			return err
		}
		if err := e.dispatcher.Dispatch(ctx, evt); err != nil {
			return err
		}
	}
}
