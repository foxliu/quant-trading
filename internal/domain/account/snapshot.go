package account

import (
	"quant-trading/internal/domain/capital"
	"quant-trading/internal/domain/portfolio"
	"time"
)

/*
Snapshot
========

账户状态快照。

设计原则:

1 不包含指针
2 不包含接口
3 完全不可变
4 可序列化
*/
type Snapshot struct {
	AccountID string

	Balance   BalanceSnapshot
	Capital   capital.Snapshot
	Portfolio portfolio.Snapshot

	RealizedPnL float64
	Timestamp   time.Time
}
