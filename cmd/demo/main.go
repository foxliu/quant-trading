package main

import (
	"fmt"
	"quant-trading/internal/application/account"
	dAccount "quant-trading/internal/domain/account"
	"quant-trading/internal/domain/market"
	"quant-trading/internal/domain/strategy"
	"time"
)

func main() {
	// 创建账户上下文
	accCfg := dAccount.Config{
		AccountID:   "demo_account",
		InitialCash: 100000.0,
	}
	accCtx := account.NewContext(accCfg)

	// 创建策略上下文
	stgCtx := strategy.NewContext()
	stgCtx.SetAccountContext(accCtx)

	// 创建行情事件
	event := market.Event{
		Type: market.EventMarket,
		Time: time.Now(),
		Data: nil,
	}
	stgCtx.SetCurrentEvent(event)
	stgCtx.SetNow(event.Time)

	fmt.Printf("账户信息: ID=%s, 现金=%.2f, 权益=%.2f\n",
		accCtx.AccountID(), accCtx.Cash(), accCtx.Equity())
	fmt.Println("策略引擎演示完成")
}
