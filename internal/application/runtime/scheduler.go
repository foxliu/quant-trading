package runtime

import (
	"context"
	"fmt"
	"sync"

	"quant-trading/internal/application/account"
	"quant-trading/internal/application/event"
	aExecution "quant-trading/internal/application/execution"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/strategy"
	"quant-trading/internal/infrastructure/logger"

	"go.uber.org/zap"
)

// Scheduler 多策略调度器（核心调度层）
type Scheduler struct {
	mu         sync.RWMutex
	runtimes   map[string]*Runtime
	accountCtx *account.Context
	bus        event.Bus
	engine     aExecution.Engine
	logger     *zap.Logger
}

// fillListener（成交监听器，保持不变）
type fillListener struct {
	accountCtx *account.Context
	logger     *zap.Logger
}

func (l *fillListener) OnExecutionEvent(ctx context.Context, evt *execution.Event) {
	if evt.Type != execution.EventOrderFilled {
		return
	}
	l.accountCtx.ApplyFill(evt.Symbol, evt.Side, evt.Price, evt.Quantity)
	l.logger.Info("成交已应用到账户",
		zap.String("orderID", evt.OrderID),
		zap.String("symbol", evt.Symbol),
		zap.Int64("qty", evt.Quantity),
	)
}

// NewScheduler 创建多策略调度器
func NewScheduler(accCtx *account.Context, bus event.Bus, engine aExecution.Engine) *Scheduler {
	s := &Scheduler{
		runtimes:   make(map[string]*Runtime),
		accountCtx: accCtx,
		bus:        bus,
		engine:     engine,
		logger:     logger.Logger.With(zap.String("module", "runtime.scheduler")),
	}

	// 注册成交监听器
	engine.RegisterListener(&fillListener{
		accountCtx: accCtx,
		logger:     s.logger,
	})

	return s
}

// RegisterStrategy 注册策略（同时订阅市场事件）
func (s *Scheduler) RegisterStrategy(stg strategy.Strategy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.runtimes[stg.Name()]; exists {
		return fmt.Errorf("策略已存在: %s", stg.Name())
	}

	rt := NewRuntime(stg, s.accountCtx, s.bus, 1024)
	s.runtimes[stg.Name()] = rt

	// 为该 Runtime 注册市场事件 Handler（使用现有 Bus 接口）
	s.bus.Subscribe(event.EventMarketPrice, (&marketEventHandler{rt: rt}).Handle)

	s.logger.Info("策略已注册并订阅市场事件", zap.String("strategy", stg.Name()))
	return nil
}

// Start / Stop（保持不变）
func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for name, rt := range s.runtimes {
		if err := rt.Start(ctx); err != nil {
			s.logger.Error("启动策略失败", zap.String("strategy", name), zap.Error(err))
			return err
		}
		s.logger.Info("策略运行时已启动", zap.String("strategy", name))
	}

	s.logger.Info("多策略调度器启动成功", zap.Int("strategy_count", len(s.runtimes)))
	return nil
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, rt := range s.runtimes {
		rt.Stop()
		s.logger.Info("策略已停止", zap.String("strategy", name))
	}
	s.logger.Info("多策略调度器已停止")
}
