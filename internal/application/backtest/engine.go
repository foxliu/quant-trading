package backtest

import (
	"context"
	"errors"
	"quant-trading/internal/application/account"
	"quant-trading/internal/application/event"
	runtime2 "quant-trading/internal/application/runtime"
	strategyengine "quant-trading/internal/application/strategy"
	dAccount "quant-trading/internal/domain/account"
	"quant-trading/internal/domain/capital"
	"quant-trading/internal/domain/portfolio"
	"quant-trading/internal/domain/strategy"
	"time"
)

/*
Engine
======

回测引擎负责:
- 历史数据回放
- 模拟撮合
- 回测时间管理
- 回测结果分析

设计原则:
- 策略代码在回测与实盘环境中完全一致
- 回测引擎只负责驱动,不改变策略逻辑
- 支持事件驱动回测
*/
type Engine struct {
	// 策略引擎
	strategyEngine *strategyengine.Engine

	// 回测时钟
	clock *Clock

	// 模拟撮合引擎
	simulator *Simulator

	// 数据源
	dataSource DataSource

	// 账户上下文
	accountCtx *account.Context

	// 回测配置
	config Config

	// 上下文
	ctx    context.Context
	cancel context.CancelFunc
}

// Config 回测配置
type Config struct {
	StartTime time.Time // 回测开始时间
	EndTime   time.Time // 回测结束时间

	InitialCash float64 // 初始资金

	Commission float64 // 手续费率
	Slippage   float64 // 滑点(百分比)
}

// NewEngine 创建回测引擎
func NewEngine(
	stg strategy.Strategy,
	dataSource DataSource,
	config Config,
) *Engine {
	// 创建回测时钟
	clock := NewClock(config.StartTime)

	capi := capital.NewCashEngine(config.InitialCash)
	acc := dAccount.Account{AccountID: "tests"}
	port := portfolio.NewSimplePortfolio()

	// 创建账户上下文
	accountCtx := account.NewContext(acc, capi, port)

	// 创建策略上下文
	stgCtx := strategy.NewContext()
	stgCtx.SetAccountContext(accountCtx)

	bus := event.NewMemoryBus()
	recorder := event.NewMemoryRecorder()
	recordingBus := event.NewRecordingBus(bus, recorder)
	// 创建策略运行时
	runtime := runtime2.NewRuntime(stg, accountCtx, recordingBus, 1024)

	// 创建策略调度器
	dispatcher := strategyengine.NewDispatcher([]*runtime2.Runtime{runtime}, nil, recordingBus)

	// 创建策略引擎
	strategyEngine := strategyengine.NewEngine(dispatcher)

	// 创建模拟撮合引擎
	simulator := NewSimulator(config.Commission, config.Slippage)

	return &Engine{
		strategyEngine: strategyEngine,
		clock:          clock,
		simulator:      simulator,
		dataSource:     dataSource,
		accountCtx:     accountCtx,
		config:         config,
	}
}

// Run 运行回测
func (e *Engine) Run(ctx context.Context) (*Result, error) {
	e.ctx, e.cancel = context.WithCancel(ctx)
	defer e.cancel()

	// 启动策略引擎
	if err := e.strategyEngine.Start(); err != nil {
		return nil, err
	}
	defer e.strategyEngine.Stop()

	// 回测主循环
	for e.clock.Now().Before(e.config.EndTime) {
		select {
		case <-e.ctx.Done():
			return nil, errors.New("backtest cancelled")
		default:
			// 获取当前时间的行情数据
			events, err := e.dataSource.GetEvents(e.clock.Now())
			if err != nil {
				return nil, err
			}

			// 处理每个行情事件
			for _, evt := range events {
				// 更新回测时钟
				e.clock.SetNow(evt.Time)

				// 将事件发送给策略引擎
				e.strategyEngine.OnMarketEvent(evt)
			}

			// 推进时钟
			e.clock.Advance(1 * time.Minute)
		}
	}

	// 生成回测报告
	result := e.generateResult()
	return result, nil
}

// Stop 停止回测
func (e *Engine) Stop() error {
	if e.cancel != nil {
		e.cancel()
	}
	return nil
}

// generateResult 生成回测结果
func (e *Engine) generateResult() *Result {
	finalCash := e.accountCtx.TotalCapital()
	finalEquity := e.accountCtx.Equity()
	realizedPnL := e.accountCtx.RealizedPnL()

	return &Result{
		StartTime:   e.config.StartTime,
		EndTime:     e.config.EndTime,
		InitialCash: e.config.InitialCash,
		FinalCash:   finalCash,
		FinalEquity: finalEquity,
		RealizedPnL: realizedPnL,
		TotalReturn: (finalEquity - e.config.InitialCash) / e.config.InitialCash,
	}
}

// Result 回测结果
type Result struct {
	StartTime   time.Time
	EndTime     time.Time
	InitialCash float64
	FinalCash   float64
	FinalEquity float64
	RealizedPnL float64
	TotalReturn float64 // 总收益率
}
