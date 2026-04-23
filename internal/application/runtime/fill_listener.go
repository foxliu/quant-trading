package runtime

import (
	"context"
	"quant-trading/internal/application/account"
	"quant-trading/internal/domain/execution"

	"go.uber.org/zap"
)

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
