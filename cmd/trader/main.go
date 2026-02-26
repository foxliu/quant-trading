package trader

import (
	"context"
	"flag"
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
	//mode := flag.String("mode", "live", "live or backtest")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 1. 创建 EventBus
	bus := event.NewMemoryBus()
	recorder := event.NewMemoryRecorder()
	recordingBus := event.NewRecordingBus(bus, recorder)

	// 2. 创建其他组件
	accountCtx := account.NewContext(dAccount.Config{})
	//riskEngine := risk.NewEngine()
	runtimes := make([]*runtime.Runtime, 0)
	registry := strategyengine.NewRegistry()
	for _, s := range registry.All() {
		runtimes = append(runtimes, runtime.NewRuntime(s, accountCtx, bus, 1024))
	}

	// 5. 创建 Dispatcher
	dispatcher := strategyengine.NewDispatcher(runtimes, nil, recordingBus)

	// 6. 启动
	if err := dispatcher.Start(ctx); err != nil {
		log.Fatal(err)
	}

	bus.Publish(&event.Envelope{
		Type:      event.EventMarketPrice,
		Timestamp: time.Now(),
		Payload:   market.Event{},
	})

	replayer := event.NewReplayer(bus, nil)
	replayer.ReplayFromSnapshot(nil, recorder.Events())

	<-ctx.Done()
}
