package risk

import (
	"context"
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/infrastructure/logger"

	"go.uber.org/zap"
)

type Coordinator struct {
	engine *Engine
	bus    event.Bus
	logger *zap.Logger
}

func NewCoordinator(engine *Engine, bus event.Bus) *Coordinator {
	return &Coordinator{
		engine: engine,
		bus:    bus,
		logger: logger.Logger.With(zap.String("module", "risk.coordinator")),
	}
}

func (c *Coordinator) Start(ctx context.Context) error {
	c.bus.Subscribe(event.EventOrderSubmitted, func(evt *event.Envelope) {
		if ord, ok := evt.Payload.(*order.Order); ok {
			result := c.engine.CheckOrder(ord)
			if result.Action.IsBlock() {
				c.logger.Warn("风控拦截订单", zap.String("message", result.Message))
			}
		}
	})
	c.logger.Info("风控协调器启动成功")
	return nil
}
