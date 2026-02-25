package event

import "sync"

type MemoryBus struct {
	mu        sync.RWMutex
	handlers  map[Type][]Handler
	nextEvent uint64
}

/*
同步调用（极其重要）

无 goroutine

Replay 时顺序 = Publish 顺序
*/

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		handlers: make(map[Type][]Handler),
	}
}

func (b *MemoryBus) Subscribe(t Type, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[t] = append(b.handlers[t], h)
}

func (b *MemoryBus) Publish(evt *Envelope) {
	b.mu.RLock()
	hs := b.handlers[evt.Type]
	b.mu.RUnlock()

	for _, h := range hs {
		h(evt)
	}
}
