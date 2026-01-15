package order

/*
Status
======

Status 表示订单在执行系统中的生命周期状态。
*/
type Status string

const (
	Pending  Status = "PENDING"  // 已生成，尚未发送
	Placed   Status = "PLACED"   // 已发送至交易所
	Filled   Status = "FILLED"   // 完全成交
	Canceled Status = "CANCELED" // 已取消
	Rejected Status = "REJECTED" // 被拒绝
)
