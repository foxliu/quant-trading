package broker

import (
	"quant-trading/internal/domain/execution"
)

type CTPReqID int

func (i CTPReqID) Value() int {
	return int(i)
}

type PendingRequest struct {
	ch chan *execution.AccountEvent
}
