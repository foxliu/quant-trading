package event

import (
	"quant-trading/internal/infrastructure/logger"
	"sync"

	"go.uber.org/zap"
)

type MemoryBus struct {
	mu       sync.RWMutex
	handlers map[Type][]Handler
	logger   *zap.Logger
}

/*
同步调用（极其重要）

无 goroutine

Replay 时顺序 = Publish 顺序
*/

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		handlers: make(map[Type][]Handler),
		logger:   logger.Logger.With(zap.String("module", "event.memory_bus")),
	}
}

func (b *MemoryBus) Subscribe(t Type, h Handler) {
	if h == nil {
		return
	}

	b.mu.Lock()
	b.handlers[t] = append(b.handlers[t], h)
	b.mu.Unlock()
	b.logger.Debug("事件已订阅",
		zap.String("type", t.String()),
		zap.String("handler", "anonymous"))
}

func (b *MemoryBus) Publish(evt *Envelope) {
	if evt == nil {
		return
	}

	b.mu.RLock()
	handlers := b.handlers[evt.Type]
	b.mu.RUnlock()

	if len(handlers) == 0 {
		return
	}

	b.logger.Debug("事件已发布",
		zap.String("type", evt.Type.String()),
		zap.String("source", evt.Source),
		zap.Int("handler_count", len(handlers)))

	// 并发安全调用所有handler
	for _, h := range handlers {
		go h(evt) // 异步执行， 防止慢 handler 阻塞总线
	}
}
