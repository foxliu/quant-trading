package position

import (
	dExecution "quant-trading/internal/domain/execution"
	"sync"
)

//多 Symbol 聚合

type Manager struct {
	mu       sync.Mutex
	contexts map[string]*Context
}

func NewManager() *Manager {
	return &Manager{
		contexts: make(map[string]*Context),
	}
}

func (m *Manager) Get(symbol string) *Context {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx, ok := m.contexts[symbol]
	if !ok {
		ctx = NewContext(symbol)
		m.contexts[symbol] = ctx
	}
	return ctx
}

func (m *Manager) OnExecutionEvent(evt *dExecution.Event, symbol string) error {
	ctx := m.Get(symbol)
	return ctx.OnExecutionEvent(evt)
}
