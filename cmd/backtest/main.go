package main

import (
	"context"
	"fmt"
	"quant-trading/internal/application/backtest"
	"quant-trading/internal/domain/instrument"
	"quant-trading/internal/infrastructure/marketdata"
	"quant-trading/internal/infrastructure/strategy/examples"
	"time"
)

func main() {
	fmt.Print("=== 量化交易系统回测示例 ===\n\n")

	// 1. 创建合约模型...
	fmt.Println("1. 创建合约模型...")
	instr := instrument.Instrument{
		Symbol:   "AAPL",
		Exchange: "NASDAQ",
		Type:     instrument.Stock,
	}
	fmt.Printf("  合约: %s (%s)\n\n", instr.Symbol, instr.Exchange)

	// 2. 使用CSVAdapter加载真实行情数据
	fmt.Println("2. 使用CSVAdapter加载真实行情数据...")
	adapter, err := marketdata.NewCSVAdapter("data/AAPL.csv", instr)
	if err != nil {
		fmt.Printf("❌ CSV 加载失败: %v\n", err)
		fmt.Println("   请确保 data/AAPL_1min.csv 文件存在并格式正确（timestamp,open,high,low,close,volume）")
		return
	}
	dataSource := adapter // *CSVAdapter 直接实现 backtest.DataSource
	fmt.Printf("   ✅ CSV 加载成功（数据源已就绪）\n\n")

	// 3. 创建策略
	fmt.Println("3. 创建双均线交叉策略...")
	strategy := examples.NewMACrossStrategy("MA_Cross_5_20", 5, 20)
	fmt.Print("   策略参数: 短期均线=5, 长期均线=20\n\n")

	// 4. 创建回测配置
	config := backtest.Config{
		StartTime:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		InitialCash: 100000.0,
		Commission:  0.0003, // 0.03% 手续费
		Slippage:    0.0001, // 0.01% 滑点
	}

	fmt.Println("4. 回测配置:")
	fmt.Printf("   开始时间: %s\n", config.StartTime.Format("2006-01-02"))
	fmt.Printf("   结束时间: %s\n", config.EndTime.Format("2006-01-02"))
	fmt.Printf("   初始资金: %.2f\n", config.InitialCash)
	fmt.Printf("   手续费率: %.4f%%\n", config.Commission*100)
	fmt.Printf("   滑点率: %.4f%%\n\n", config.Slippage*100)

	// 5. 创建回测引擎
	fmt.Println("5. 创建回测引擎...")
	engine := backtest.NewEngine(strategy, dataSource, config)
	fmt.Print("   回测引擎已就绪（使用 CSV 数据源）\n\n")

	// 6. 运行回测
	fmt.Println("6. 开始回测...")
	result, err := engine.Run(context.Background())
	if err != nil {
		fmt.Printf("回测失败: %v\n", err)
		return
	}
	fmt.Printf("   回测完成!\n\n")

	// 7. 输出结果
	fmt.Println("=== 回测结果 ===")
	fmt.Printf("初始资金: %.2f\n", result.InitialCash)
	fmt.Printf("最终现金: %.2f\n", result.FinalCash)
	fmt.Printf("最终权益: %.2f\n", result.FinalEquity)
	fmt.Printf("已实现盈亏: %.2f\n", result.RealizedPnL)
	fmt.Printf("总收益率: %.2f%%\n", result.TotalReturn*100)
	fmt.Println("\n回测完成! 使用 CSV 数据源成功 ✅")
}
