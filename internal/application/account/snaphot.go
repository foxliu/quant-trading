package account

import (
	"quant-trading/internal/application/position"
	"quant-trading/internal/domain/account"
)

/*
唯一对外出口
*/

type Snapshot struct {
	AccountID string

	Balance   account.Balance
	Positions map[string]*position.Snapshot
}

func (c *Context) Snapshot() *Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()

	posCopy := make(map[string]*position.Snapshot)
	for k, v := range c.positions {
		posCopy[k] = v
	}

	return &Snapshot{
		AccountID: c.account.AccountID,
		Balance:   c.balance,
		Positions: posCopy,
	}
}
