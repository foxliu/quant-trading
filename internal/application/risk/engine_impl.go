package risk

import (
	"context"
	"quant-trading/internal/domain/strategy"
	"sync"
)

type engine struct {
	ctx    context.Context
	cancel context.CancelFunc

	signalCh chan strategy.Signal
	wg       sync.WaitGroup

	ctxProvider ContextProvider
}

func NewEngine(ctxProvider ContextProvider, buffer int) Engine {
	return &engine{
		signalCh:    make(chan strategy.Signal, buffer),
		ctxProvider: ctxProvider,
	}
}

func (e *engine) Start(parent context.Context) error {
	e.ctx, e.cancel = context.WithCancel(parent)

	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.run()
	}()
	return nil
}

func (e *engine) Stop() error {
	e.cancel()
	e.wg.Wait()
	return nil
}

func (e *engine) Consume(signal strategy.Signal) {
	select {
	case e.signalCh <- signal:
	default:
		// 丢弃， 保证不阻塞上游
	}
}
