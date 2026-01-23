package account

/*
Balance
=======

Balance 描述账户的资金状态
*/
type Balance struct {
	Cash        float64 // 可用现金
	Frozen      float64 // 冻结资金
	Equity      float64 // 总权益
	RealizedPnL float64 // 已实现盈亏
}
