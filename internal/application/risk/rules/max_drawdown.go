package rules

import (
	"quant-trading/internal/application/risk"
	risk2 "quant-trading/internal/domain/risk"
	"time"
)

/*
最大回撤（Floating Loss)
*/

type MaxDrawDown struct {
	Limit float64
}

func (r *MaxDrawDown) Name() string {
	return "MaxDrawdown"
}

func (r *MaxDrawDown) Evaluate(ctx *risk.Context) *risk.Result {
	ctx.Mu.Lock()
	defer ctx.Mu.Unlock()

	if ctx.PnL.Unrealized < -r.Limit {
		return &risk.Result{
			RuleName: r.Name(),
			Action:   risk2.ActionForceClose,
			Reason:   "floating loss exceeds limit",
			Time:     time.Now(),
		}
	}
	return nil
}
