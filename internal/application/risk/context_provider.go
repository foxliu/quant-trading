package risk

import "quant-trading/internal/domain/strategy"

/*
ContextProvider
===============

用于根据 Signal 获取 Risk Context。

这是 Risk Engine 与 Account / Runtime 解耦的关键。
*/
type ContextProvider interface {
	ContextFor(signal strategy.Signal) Context
}
