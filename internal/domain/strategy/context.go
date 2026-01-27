package strategy

import (
	"quant-trading/internal/domain/market"
	"time"
)

/*
Context
=======

Context 是 Strategy 运行时的上下文对象。

设计原则：
1. 每个 Strategy 实例独占一个 Context
2. Context 可读、可写（仅影响该策略）
3. Context 不允许产生副作用（如下单）
*/
type Context interface {

	// ========= 基础时间 =========

	Now() time.Time
	SetNow(t time.Time)

	// ========= 行情快照 =========

	CurrentEvent() market.Event
	SetCurrentEvent(event market.Event)

	// ========= 策略私有状态 =========

	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	MustGet(key string) interface{}
	Delete(key string)

	// ========= 参数（启动时注入） =========

	Params() map[string]interface{}
	SetParams(params map[string]interface{})

	// ========= Account =========

	Account() AccountReader
	SetAccountContext(ctx AccountReader)
}
