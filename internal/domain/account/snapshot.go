package account

import (
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
	AccountID  string // 账户唯一标识
	TradingDay string // CTP 返回的交易日（格式 YYYYMMDD）
	Balance    BalanceSnapshot
	UpdateTime time.Time // 快照更新时间（UTC）
	Version    int64     // 版本号（原子递增，用于并发防覆盖）
	// 可扩展字段（未来风控/组合需要）
	// MarginRatio float64      // 维持保证金比例（可选）
	// RiskFlags   []RiskFlag    // 预留风控标记位
}
