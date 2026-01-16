package account

import (
	"quant-trading/internal/application/position"
	"quant-trading/internal/domain/account"
	"sync"
)

type Context struct {
	mu sync.Mutex

	account account.Account

	balance account.Balance

	positions map[string]*position.Snapshot
}

func NewContext(acc account.Account, initialCash float64) *Context {
	return &Context{
		account: acc,
		balance: account.Balance{
			Cash:   initialCash,
			Equity: initialCash,
		},
		positions: make(map[string]*position.Snapshot),
	}
}
