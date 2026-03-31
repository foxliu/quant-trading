package risk

import "quant-trading/internal/domain/order"

// Context 风控上下文（Strategy 只能通过此上下文与风控交互）
type Context struct {
	engine RiskEngine
}

func NewContext(engine RiskEngine) *Context {
	return &Context{engine: engine}
}

// CheckOrder 下单前风控检查
func (c *Context) CheckOrder(ord *order.Order) CheckResult {
	return c.engine.CheckOrder(ord)
}

// CheckPosition 持仓风控检查
func (c *Context) CheckPosition() CheckResult {
	return c.engine.CheckPosition()
}

// GetStatus 获取当前风控检查
func (c *Context) GetStatus() Status {
	return c.engine.GetStatus()
}

// EmergencyStop 主动触发紧急停止
func (c *Context) EmergencyStop(reason string) {
	c.engine.EmergencyStop(reason)
}
