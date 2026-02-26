package risk

import (
	"context"
)

type Coordinator interface {
	Start(ctx context.Context) error
	Stop() error

	OnRiskResult(res *Result)
}
