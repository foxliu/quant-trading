package pnl

import "context"

type Engine interface {
	Start(ctx context.Context) error
	Stop() error
}
