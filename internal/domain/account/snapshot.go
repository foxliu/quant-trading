package account

import (
	"quant-trading/internal/domain/capital"
	"quant-trading/internal/domain/common"
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
	common.Model
	AccountID   AccountID          `gorm:"column:account_id;index;not null"` // 账户唯一标识
	TradingDay  string             `gorm:"column:trading_day"`               // CTP 返回的交易日（格式 YYYYMMDD）
	Balance     BalanceSnapshot    `gorm:"column:balance;type:JSONB"`
	Capital     capital.Snapshot   `gorm:"column:capital;type:JSONB"`
	Portfolio   portfolio.Snapshot `gorm:"column:portfolio;type:JSONB"`
	RealizedPnL float64            `gorm:"column:realized_pnl"`
	Timestamp   time.Time          `gorm:"index"`
}

func (s *Snapshot) TableName() string {
	return "t_account_snapshot"
}
