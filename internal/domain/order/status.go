package order

/*
Status
======

Status 表示订单在执行系统中的生命周期状态。
*/
type Status int

const (
	StatusPending  Status = iota // 已生成，尚未发送
	StatusPlaced                 // 已发送至交易所
	StatusFilled                 // 完全成交
	StatusCanceled               // 已取消
	StatusRejected               // 被拒绝
)
