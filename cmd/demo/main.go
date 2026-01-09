package main

import (
	"context"
	"fmt"
	"quant-trading/internal/application/strategy"
	"quant-trading/internal/domain"
	"time"
)

func main() {
	bar := domain.MarketBar{
		Instrument: "AAPL",
		Time:       time.Now(),
		Close:      120,
	}

	sc := strategy.NewContext(bar.Time)
	stg := strategy.NewDemoStrategy("demo")

	engine := strategy.NewEngine(stg)

	signals := engine.RunOnBar(context.Background(), sc, bar)

	for _, s := range signals {
		fmt.Printf("Signal: %+v\n", *s)
	}
}
