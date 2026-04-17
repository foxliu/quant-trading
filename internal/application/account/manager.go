package account

import (
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/account"
	"quant-trading/internal/domain/capital"
	"quant-trading/internal/domain/portfolio"
	"quant-trading/internal/infrastructure/logger"
	"sync"

	"go.uber.org/zap"
)

type Manager struct {
	service  *Service
	contexts map[account.AccountID]*Context
	mu       sync.RWMutex
	log      *zap.Logger
}

func NewManager(service *Service) *Manager {
	return &Manager{
		service:  service,
		contexts: make(map[account.AccountID]*Context),
		log:      logger.Logger.With(zap.String("module", "account.manager")),
	}
}

func (m *Manager) GetContext(id account.AccountID, cap capital.Engine, port portfolio.Engine) (*Context, error) {
	m.mu.RLock()
	if ctx, ok := m.contexts[id]; ok {
		m.mu.RUnlock()
		return ctx, nil
	}
	m.mu.RUnlock()

	ctx, err := m.service.NewContext(id, cap, port)
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
	m.contexts[id] = ctx
	m.mu.Unlock()
	return ctx, nil
}

func (m *Manager) StartPersister() {
	m.service.bus.Subscribe(event.EventAccountBalanceChanged, func(evt *event.Envelope) {
		snap, ok := evt.Payload.(account.Snapshot)
		if !ok {
			m.log.Error("无效的evt.Payload")
			return
		}
		m.service.repo.SaveSnapshot(&snap)
	})
}
