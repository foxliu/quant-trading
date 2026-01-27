package account

import (
	"quant-trading/internal/domain/account"
	"time"
)

/*
唯一对外出口
*/

type Snapshot struct {
	Balance account.Balance
	At      time.Time
}

func (s *Snapshot) Name() string {
	return "account"
}

func (s *Snapshot) Timestamp() time.Time {
	return s.At
}
