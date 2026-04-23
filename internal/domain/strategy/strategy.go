package strategy

import (
	"quant-trading/internal/domain/market"
)

type Strategy interface {
	Name() string
	OnInit(ctx Context) error
	OnMarketEvent(ctx Context, evt market.Event) ([]Signal, error)
	OnStop(ctx Context) error
}
