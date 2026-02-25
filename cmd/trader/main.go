package trader

import (
	"context"
	"log"
	"quant-trading/internal/application/account"
	"quant-trading/internal/application/event"
	"quant-trading/internal/application/runtime"
	strategyengine "quant-trading/internal/application/strategy"
	dAccount "quant-trading/internal/domain/account"
	"quant-trading/internal/domain/market"
	"time"
)

func main() {
	// 1. 创建 EventBus
	bus := event.NewMemoryBus()

	// 2. 创建 Risk Engine
	//riskEngine := risk.NewEngine()

	// 3. 关联用户
	accountCtx := account.NewContext(dAccount.Config{})

	// 4. 创建所有 Runtime
	runtimes := make([]*runtime.Runtime, 0)
	registry := strategyengine.NewRegistry()
	for _, s := range registry.All() {
		runtimes = append(runtimes, runtime.NewRuntime(s, accountCtx, 1024))
	}

	// 5. 创建 Dispatcher
	dispatcher := strategyengine.NewDispatcher(runtimes, nil, bus)

	// 6. 启动
	if err := dispatcher.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	bus.Publish(&event.Envelope{
		Type:      event.EventMarketPrice,
		Timestamp: time.Now(),
		Payload:   market.Event{},
	})
}
