package rules

import (
	"quant-trading/internal/application/risk"
	risk2 "quant-trading/internal/domain/risk"
	"quant-trading/pkg/utils"
	"time"
)

/*
最大仓位限制
*/

type MaxPosition struct {
	MaxQty int64
}

func (r *MaxPosition) Name() string {
	return "MaxPosition"
}

func (r *MaxPosition) Evaluate(ctx *risk.Context) *risk.Result {
	ctx.Mu.Lock()
	defer ctx.Mu.Unlock()

	if ctx.Position.Pos != nil && utils.Abs(ctx.Position.Pos.Qty) > r.MaxQty {
		return &risk.Result{
			RuleName: r.Name(),
			Action:   risk2.ActionRejectOrder,
			Reason:   "position size exceeds limit",
			Time:     time.Now(),
		}
	}
	return nil
}
