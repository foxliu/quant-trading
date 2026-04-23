package account

import (
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/account"
	"quant-trading/internal/infrastructure/logger"
	"time"

	"go.uber.org/zap"
)

// Persister 负责监听 AccountBalanceChangedEvent 并持久化 Snapshot
// 建议在main.go中启动一次： persister := NewPersister(repo, bus); persister.Start()
type Persister struct {
	repo account.Repository
	bus  event.Bus
	log  *zap.Logger
}

func NewPersister(repo account.Repository, bus event.Bus) *Persister {
	return &Persister{
		repo: repo,
		bus:  bus,
		log:  logger.Logger.With(zap.String("module", "account.persister")),
	}
}

// Start 启动事件监听（在main.go中调用一次即可）
func (p *Persister) Start() {
	p.bus.Subscribe(event.EventAccountBalanceChanged, p.handleBalanceChanged)
	p.log.Info("[AccountPersister] 已启动，监听 AccountBalanceChangedEvent")
}

// handleBalanceChanged 是事件处理函数
func (p *Persister) handleBalanceChanged(evt *event.Envelope) {
	changedEvt, ok := evt.Payload.(account.AccountBalanceChangedEvent)
	if !ok {
		p.log.Error("[AccountPersister] 事件类型断言失败", zap.Any("payload", evt.Payload))
		return
	}

	// 异步保存防止阻塞EventBus
	go func() {
		start := time.Now()

		err := p.repo.SaveSnapshot(&changedEvt.Snapshot)
		if err != nil {
			p.log.Error("[AccountPersister] 保存Snapshot失败", zap.Error(err), zap.Any("snapshot", &changedEvt.Snapshot))
			return
		}
		p.log.Info("[AccountPersister] Snapshot 已持久化",
			zap.String("AccountID", changedEvt.AccountID.String()),
			zap.Int32("耗时ms", int32(time.Since(start).Milliseconds())),
		)
	}()
}

// Stop 可选：如果需要优雅停止，可在这里取消订阅（当前 EventBus 未提供 Unsubscribe 可省略）
func (p *Persister) Stop() {
	// todo: 未来 EventBus 支持取消订阅时实现
	p.log.Info("[AccountPersister] 已停止")
}
