package strategy

/*
AccountReader
=============

AccountReader 是策略上下文中访问账户信息的只读接口。

设计目的:
- 打破 domain/strategy 对 application/account 的依赖
- 策略只需要读取账户快照,不需要直接操作账户上下文
- 符合依赖倒置原则(DIP)
*/
type AccountReader interface {
	// AccountID 返回账户ID
	AccountID() string

	// Cash 返回可用现金
	Cash() float64

	// Equity 返回账户权益
	Equity() float64

	// RealizedPnL 返回已实现盈亏
	RealizedPnL() float64

	// UnrealizedPnL 返回未实现盈亏
	UnrealizedPnL() float64
}
