package main

import (
	"context"
	"fmt"
	"quant-trading/internal/application/backtest"
	"quant-trading/internal/domain/instrument"
	"quant-trading/internal/domain/market"
	"quant-trading/internal/infrastructure/strategy/examples"
	"time"
)

func main() {
	fmt.Println("=== 量化交易系统回测示例 ===\n")

	// 1. 创建测试数据
	fmt.Println("1. 生成测试行情数据...")
	events := generateTestData()
	fmt.Printf("   生成了 %d 条行情数据\n\n", len(events))

	// 2. 创建数据源
	dataSource := backtest.NewMemoryDataSource(events)

	// 3. 创建策略
	fmt.Println("2. 创建双均线交叉策略...")
	strategy := examples.NewMACrossStrategy("MA_Cross_5_20", 5, 20)
	fmt.Println("   策略参数: 短期均线=5, 长期均线=20\n")

	// 4. 配置回测
	config := backtest.Config{
		StartTime:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		InitialCash: 100000.0,
		Commission:  0.0003, // 0.03% 手续费
		Slippage:    0.0001, // 0.01% 滑点
	}

	fmt.Println("3. 回测配置:")
	fmt.Printf("   开始时间: %s\n", config.StartTime.Format("2006-01-02"))
	fmt.Printf("   结束时间: %s\n", config.EndTime.Format("2006-01-02"))
	fmt.Printf("   初始资金: %.2f\n", config.InitialCash)
	fmt.Printf("   手续费率: %.4f%%\n", config.Commission*100)
	fmt.Printf("   滑点率: %.4f%%\n\n", config.Slippage*100)

	// 5. 创建回测引擎
	fmt.Println("4. 创建回测引擎...")
	engine := backtest.NewEngine(strategy, dataSource, config)
	fmt.Println("   回测引擎已就绪\n")

	// 6. 运行回测
	fmt.Println("5. 开始回测...")
	result, err := engine.Run(context.Background())
	if err != nil {
		fmt.Printf("回测失败: %v\n", err)
		return
	}
	fmt.Println("   回测完成!\n")

	// 7. 输出结果
	fmt.Println("=== 回测结果 ===")
	fmt.Printf("初始资金: %.2f\n", result.InitialCash)
	fmt.Printf("最终现金: %.2f\n", result.FinalCash)
	fmt.Printf("最终权益: %.2f\n", result.FinalEquity)
	fmt.Printf("已实现盈亏: %.2f\n", result.RealizedPnL)
	fmt.Printf("总收益率: %.2f%%\n", result.TotalReturn*100)
	fmt.Println("\n回测完成!")
}

// generateTestData 生成测试行情数据
func generateTestData() []market.Event {
	events := make([]market.Event, 0)

	// 生成30天的模拟行情数据
	startTime := time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC)
	basePrice := 100.0

	for day := 0; day < 30; day++ {
		for minute := 0; minute < 240; minute++ { // 每天240分钟交易时间
			t := startTime.Add(time.Duration(day*24+minute) * time.Minute)

			// 简单的正弦波模拟价格波动
			price := basePrice + 10*float64(day)/30.0 + 2*float64(minute%60)/60.0

			bar := market.Bar{
				Instrument: instrument.Instrument{
					Symbol:   "AAPL",
					Exchange: "NASDAQ",
					Type:     instrument.Stock,
				},
				Time:   t,
				Open:   price - 0.1,
				High:   price + 0.2,
				Low:    price - 0.2,
				Close:  price,
				Volume: 1000000,
			}

			event := market.Event{
				Type: market.EventMarket,
				Time: t,
				Data: bar,
			}

			events = append(events, event)
		}
	}

	return events
}
