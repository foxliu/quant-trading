package risk

import (
	"context"
	"sync"
	"time"
)

type engine struct {
	ctx *Context

	rules []Rule

	resultCh chan *Result
	stopCh   chan struct{}

	mu sync.Mutex
}

func NewEngine(ctx *Context, rules ...Rule) Engine {
	return &engine{
		ctx:      ctx,
		rules:    rules,
		resultCh: make(chan *Result, 64),
		stopCh:   make(chan struct{}),
	}
}

func (e *engine) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-e.stopCh:
				return
			default:
				e.Evaluate()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	return nil
}

func (e *engine) Stop() error {
	close(e.stopCh)
	return nil
}

func (e *engine) Evaluate() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, rule := range e.rules {
		if res := rule.Evaluate(e.ctx); res != nil {
			e.resultCh <- res
		}
	}
}

func (e *engine) Results() <-chan *Result {
	return e.resultCh
}
