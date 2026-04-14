package account

/*
BalanceSnapshot
===============

余额快照
*/
type BalanceSnapshot struct {
	PreBalance    float64 //  上次结算准备金
	Available     float64 //  可用资金
	Commission    float64 //  手续费
	UnrealizedPnL float64 //  未实现盈亏
	RealizedPnL   float64 //  已实现盈亏
	Frozen        float64 //  冻结的保证金
}
