package account

/*
Balance
=======

Balance 描述账户的资金状态
*/
type Balance struct {
	Equity      float64 // 总权益
	RealizedPnL float64 // 已实现盈亏
}
