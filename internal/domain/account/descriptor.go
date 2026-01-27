package account

/*
账户是谁, 是配置，不是状态
*/

type Descriptor struct {
	AccountID string
	Exchange  string
	Currency  string
}
