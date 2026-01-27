package account

import (
	"errors"
	"quant-trading/internal/application/snapshot"
	"time"
)

func (c *Context) Take() snapshot.Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()

	b := c.balance
	return &Snapshot{
		Balance: b,
		At:      time.Now(),
	}
}

func (c *Context) Restore(s snapshot.Snapshot) error {
	as, ok := s.(*Snapshot)
	if !ok {
		return errors.New("invalid account snapshot")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.balance = as.Balance
	return nil
}
