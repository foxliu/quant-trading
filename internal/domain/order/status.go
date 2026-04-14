package order

/*
OrderStatus
======

OrderStatus 表示订单在执行系统中的生命周期状态。
*/
type OrderStatus int

const (
	OrderStatusAllTraded             OrderStatus = iota // 全部成交
	OrderStatusPartTradedQueueing                       // 部分成交还在队列中
	OrderStatusPartTradedNotQueueing                    // 部分成交不在队列中
	OrderStatusNoTradeQueueing                          // 未成交还在队列中
	OrderStatusNoTradeNotQueueing                       // 未成交不在队列中
	OrderStatusCanceled                                 // 撤单
	OrderStatusUnknown                                  // 未知
	OrderStatusNotTouched                               // 尚未触发
	OrderStatusTouched                                  // 已触发
)
