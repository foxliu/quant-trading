package risk

/*
Context
=======

Risk Context 表示一组风控规则的集合。
*/
type Context struct {
	MaxPositionQty float64 // 单标的最大绝对仓位
	MaxOrderQty    float64 // 单笔最大下单量
}

func NewContext() *Context {
	return &Context{
		MaxPositionQty: 0, // 0 = 不限制
		MaxOrderQty:    0, // 0 = 不限制
	}
}
