package account

/*
Account
=======

账户聚合根。

这里只保留基础身份信息。
*/

type Account struct {
	AccountID string
	//mu        sync.Mutex
	//
	//id        string
	//name      string
	//capital   *Capital
	//portfolio *Portfolio
}

/*
func NewAccount(id, name string, initial float64) *Account {
	return &Account{
		id:        id,
		name:      name,
		capital:   newCapital(initial),
		portfolio: newPortfolio(),
	}
}

func (a *Account) ID() string {
	return a.id
}

func (a *Account) Name() string {
	return a.name
}

// =========================
// 核心交易入口（唯一写入口）
// =========================

func (a *Account) ApplyFill(fill Fill) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := a.portfolio.applyFill(fill); err != nil {
		return err
	}

	if err := a.capital.applyFill(fill); err != nil {
		// 若资金失败，回滚持仓
		a.portfolio.rollbackFill(fill)
		return err
	}
	return nil
}

// =========================
// 只读接口
// =========================

func (a *Account) Snapshot(market map[string]float64) Snapshot {
	a.mu.Lock()
	defer a.mu.Unlock()

	positions, totalMarketValue, totalUnrealized := a.portfolio.snapshot(market)

	var totalRealized float64
	for _, p := range positions {
		totalRealized += p.RealizedPnL
	}

	equity := a.capital.available + totalMarketValue

	return Snapshot{
		AccountID:     a.AccountID,
		Equity:        equity,
		RealizedPnL:   totalRealized,
		UnrealizedPnL: totalUnrealized,
		Positions:     positions,
	}
}

*/
