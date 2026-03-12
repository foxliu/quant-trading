package capital

/* 资本模型抽象
不再让 AccountContext 来处理 cash - frozen
而是处理 capital engine
这样未来可以切换：
Cash Account  现金账户  A股交易使用的就是现金账户
Margin Account   保证金账户  期货期权使用是保证金账户
Portfolio Margin   投资组合保证金
Futures Cross Margin    期货交叉保证金
*/

import (
	"quant-trading/internal/domain/execution"
)

type Engine interface {
	Freeze(orderID string, symbol string, price float64, qty float64, side execution.Side) error
	Commit(orderID string, amount float64) error
	Release(orderID string) error

	Available() float64
	Frozen() float64
	Total() float64

	Snapshot() Snapshot
	Restore(snapshot Snapshot)
}
