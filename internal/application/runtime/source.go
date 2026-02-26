package runtime

import (
	"context"
	"quant-trading/internal/application/event"
)

type EventSource interface {
	Next(ctx context.Context) (*event.Envelope, bool)
}
