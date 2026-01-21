package pnl

import "time"

/*
Snapshot
=====

# 是只读快照

可以被：

* Risk

* Monitor

* UI

* Recorder

直接消费
*/
type Snapshot struct {
	Symbol string

	// === 仓位快照 ===
	Qty      int64
	AvePrice float64

	// === 盈亏 ===
	Realized   float64
	Unrealized float64

	// === 估值 ===
	MarkPrice float64

	UpdateAt time.Time
}
