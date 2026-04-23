# 量化交易系统 (Quant Trading System)

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-active-success.svg)]()

一个基于Go语言实现的量化交易系统,采用DDD(领域驱动设计) + Engine架构模式。

## 🎯 核心特性

根据[平台设计宪章](docs/平台设计宪章.md),系统实现了以下四大核心能力:

- ✅ **多策略同时运行,互相隔离** - 通过Dispatcher + Runtime架构实现
- ✅ **多账户统一调度与风险约束** - 通过AccountScheduler实现
- ✅ **跨资产类型的一致策略体验** - 支持股票/期货/期权
- ✅ **同一策略代码可运行于回测与实盘** - 统一策略接口

## 📁 项目结构

```
quant-trading/
├── cmd/                    # 命令行程序
│   ├── demo/              # 演示程序
│   ├── backtest/          # 回测程序
│   └── trader/            # 实盘交易程序
├── internal/              # 内部代码
│   ├── domain/            # 领域层(纯业务逻辑)
│   │   ├── common/        # 公共枚举
│   │   ├── market/        # 行情模型
│   │   ├── instrument/    # 资产模型
│   │   ├── strategy/      # 策略接口
│   │   ├── order/         # 订单模型
│   │   ├── trade/         # 成交模型
│   │   ├── account/       # 账户模型
│   │   ├── execution/     # 执行事件
│   │   └── risk/          # 风控模型
│   ├── application/       # 应用层(业务编排)
│   │   ├── strategy/      # 策略引擎
│   │   ├── account/       # 账户管理
│   │   ├── position/      # 仓位管理
│   │   ├── risk/          # 风控引擎
│   │   ├── execution/     # 执行引擎
│   │   ├── backtest/      # 回测引擎
│   │   ├── instrument/    # 资产管理
│   │   ├── market/        # 行情管理
│   │   └── event/         # 事件总线
│   └── infrastructure/    # 基础设施层
│       └── strategy/      # 策略实现
│           └── examples/  # 示例策略
├── docs/                  # 文档
│   ├── 平台设计宪章.md
│   └── ...
├── README.md
├── analysis_report.md     # 系统分析报告
└── COMPLETION_REPORT.md   # 完成报告
```

## 🚀 快速开始

### 环境要求

- Go 1.24+
- Linux/macOS/Windows

### 安装

```bash
# 克隆项目
git clone <repository-url>
cd quant-trading

# 编译项目
go build ./...
```

### 运行Demo

```bash
# 运行演示程序
go run cmd/demo/main.go
```

输出:
```
账户信息: ID=demo_account, 现金=100000.00, 权益=100000.00
策略引擎演示完成
```

### 运行回测

```bash
# 编译回测程序
go build -o bin/backtest cmd/backtest/main.go

# 运行回测
./bin/backtest
```

输出:
```
=== 量化交易系统回测示例 ===
1. 生成测试行情数据...
   生成了 7200 条行情数据
2. 创建双均线交叉策略...
   策略参数: 短期均线=5, 长期均线=20
3. 回测配置:
   开始时间: 2024-01-01
   结束时间: 2024-01-31
   初始资金: 100000.00
   手续费率: 0.0300%
   滑点率: 0.0100%
4. 创建回测引擎...
   回测引擎已就绪
5. 开始回测...
   回测完成!
=== 回测结果 ===
初始资金: 100000.00
最终现金: 100000.00
最终权益: 100000.00
已实现盈亏: 0.00
总收益率: 0.00%
回测完成!
```

## 📚 核心概念

### 1. 策略(Strategy)

策略是交易逻辑的核心,实现`Strategy`接口:

```go
type Strategy interface {
    Name() string
    OnInit(ctx Context) error
    OnMarketEvent(ctx Context, evt market.Event) ([]Signal, error)
    OnStop(ctx Context) error
}
```

### 2. 信号(Signal)

策略产生的交易意图:

```go
type Signal struct {
    StrategyID string
    Symbol     string
    Intent     PositionIntent  // LONG/SHORT/FLAT
    TargetQty  float64
    Price      float64
}
```

### 3. 账户(Account)

管理资金和持仓:

```go
type Context struct {
    cfg      account.Config
    balance  account.Balance
    updateAt time.Time
}
```

### 4. 资产适配器(Adapter)

处理不同资产类型的差异:

```go
type Adapter interface {
    GetType() instrument.Type
    ValidateOrder(ord *order.Order) error
    CalculateMargin(qty int64, price float64) float64
    CalculateValue(qty int64, price float64) float64
}
```

## 💡 示例策略

### 双均线交叉策略

```go
strategy := examples.NewMACrossStrategy("MA_Cross_5_20", 5, 20)
```

**策略逻辑:**
- 短期均线上穿长期均线 → 买入
- 短期均线下穿长期均线 → 卖出

### 突破策略

```go
strategy := examples.NewBreakoutStrategy("Breakout_20", 20, 100)
```

**策略逻辑:**
- 价格突破N日最高价 → 买入
- 价格跌破N日最低价 → 卖出

## 🏗️ 架构设计

### 分层架构

```
┌─────────────────────────────────────┐
│         Interfaces Layer            │  对外接口(CLI/gRPC/HTTP)
├─────────────────────────────────────┤
│       Application Layer             │  应用层(Engine/Coordinator)
├─────────────────────────────────────┤
│         Domain Layer                │  领域层(纯业务逻辑)
├─────────────────────────────────────┤
│     Infrastructure Layer            │  基础设施(Broker/Storage)
└─────────────────────────────────────┘
```

### 交易主链路

```
Market Event → Strategy Engine → Signal → Planner → Order
    ↓
Risk Engine → Execution Engine → Broker → Execution Event
    ↓
Position Context → Account Context
```

### 多策略隔离

```
Dispatcher
  ├── Runtime 1 (Strategy A) → Channel 1
  ├── Runtime 2 (Strategy B) → Channel 2
  └── Runtime 3 (Strategy C) → Channel 3
```

## 🔧 开发指南

### 创建自定义策略

```go
package mystrategy

import (
    "quant-trading/internal/domain/market"
    "quant-trading/internal/domain/strategy"
)

type MyStrategy struct {
    name string
}

func (s *MyStrategy) Name() string {
    return s.name
}

func (s *MyStrategy) OnInit(ctx strategy.Context) error {
    // 初始化逻辑
    return nil
}

func (s *MyStrategy) OnMarketEvent(ctx strategy.Context, evt market.Event) ([]strategy.Signal, error) {
    signals := make([]strategy.Signal, 0)
    
    // 策略逻辑
    signals = append(signals, strategy.Signal{
        StrategyID: s.name,
        Symbol:     "AAPL",
        Intent:     strategy.IntentLong,
        TargetQty:  100,
        Price:      150.0,
    })
    
    return signals, nil
}

func (s *MyStrategy) OnStop(ctx strategy.Context) error {
    return nil
}
```

### 运行回测

```go
// 1. 创建策略
strategy := mystrategy.NewMyStrategy("my_strategy")

// 2. 准备数据
dataSource := backtest.NewMemoryDataSource(events)

// 3. 配置回测
config := backtest.Config{
    StartTime:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
    EndTime:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
    InitialCash: 100000.0,
    Commission:  0.0003,
    Slippage:    0.0001,
}

// 4. 运行回测
engine := backtest.NewEngine(strategy, dataSource, config)
result, err := engine.Run(context.Background())
```

## 📖 文档

- [平台设计宪章](docs/平台设计宪章.md) - 系统设计原则
- [系统分析报告](analysis_report.md) - 现状分析与差距评估
- [完成报告](COMPLETION_REPORT.md) - 功能完成情况

## 🛣️ Roadmap

### v1.0 (已完成) ✅
- [x] 多策略隔离
- [x] 多账户调度
- [x] 跨资产支持(股票/期货/期权)
- [x] 回测引擎
- [x] 示例策略(MA Cross, Breakout)

### v1.1 (计划中)
- [ ] 完善回测指标(夏普比率/最大回撤)
- [ ] CSV数据源支持
- [ ] 更多风控规则
- [ ] 策略性能分析

### v2.0 (未来)
- [ ] 实盘交易支持
- [ ] 券商API对接
- [ ] Web监控界面
- [ ] 分布式部署

## 🤝 贡献

欢迎提交Issue和Pull Request!

## 📄 许可证

MIT License

## 📧 联系方式

- 项目主页: <repository-url>
- 问题反馈: <issues-url>

---

**Built with ❤️ using Go**
