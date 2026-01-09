package strategy

/*
最小定义
*/

import "time"

// Context 是策略运行时上下文
// 设计原则：
// 1. 对策略只读
// 2. 不允许直接下单
// 3. 不绑定交易所实现
type Context struct {
	now time.Time
}

func NewContext(now time.Time) Context {
	return Context{
		now: now,
	}
}

func (c *Context) Now() time.Time {
	return c.now
}
