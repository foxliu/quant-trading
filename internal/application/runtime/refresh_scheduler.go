package runtime

import (
	"context"
	"quant-trading/internal/application/account"
	"quant-trading/internal/application/event"
	dAccount "quant-trading/internal/domain/account"
	"quant-trading/internal/infrastructure/logger"
	"time"

	"go.uber.org/zap"
)

type RefreshScheduler struct {
	ctx    *account.Context
	ticker *time.Ticker
	cancel context.CancelFunc
	bus    event.Bus
	log    *zap.Logger
}

func NewRefreshScheduler(ctx *account.Context, bus event.Bus) *RefreshScheduler {
	return &RefreshScheduler{
		ctx: ctx,
		bus: bus,
		log: logger.Logger.With(zap.String("module", "runtime.refresh_scheduler")),
	}
}

func (s *RefreshScheduler) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	s.ticker = time.NewTicker(500 * time.Millisecond) // 500 ms

	go func() {
		for {
			select {
			case <-s.ticker.C:
				snap := s.ctx.Snapshot()
				s.log.Info("[RefreshScheduler] running",
					zap.Float64("equity", s.ctx.Equity()),
					zap.String("account_id", snap.AccountID.String()),
					zap.Float64("realizePnl", s.ctx.RealizedPnL()),
				)

				// 每5秒强制发布一次快照事件
				if time.Now().Unix()%5 == 0 {
					s.bus.Publish(&event.Envelope{
						Type:      event.EventAccountBalanceChanged,
						Source:    "refresh.scheduler",
						Timestamp: time.Now(),
						Payload: dAccount.AccountBalanceChangedEvent{
							AccountID: snap.AccountID,
							Snapshot:  snap,
						},
					})
				}
			case <-ctx.Done():
				s.log.Info("[RefreshScheduler] stopped")
				return
			}
		}
	}()
}

func (s *RefreshScheduler) Stop() {
	s.ticker.Stop()
	if s.cancel != nil {
		s.cancel()
	}
	s.log.Info("[RefreshScheduler] stopped")
}
