## 1️⃣ 我们当前讨论的“层级”

**只讨论：单进程量化交易系统内核（非分布式）**

* 中低频（秒级）
* 股票 / 期权 / 期货
* 实盘 + 回测复用
* Go 实现
* DDD + Engine 架构

不讨论：

* 微服务
* 分布式撮合
* 高频（毫秒 / 微秒）

---

## 2️⃣ 三个核心概念（后面不再变化）

| 概念                 | 定义                   |
| ------------------ | -------------------- |
| **Domain**         | 业务语言与规则，不关心“怎么跑”     |
| **Engine**         | 运行机制（并发 / 生命周期 / 调度） |
| **Infrastructure** | 外部世界（行情 / 券商 / 存储）   |

👉 **Strategy 是 Domain**
👉 **Strategy Engine 是 Application / Engine**


# 二、最终【唯一权威】工程目录结构

下面这份结构，**覆盖你之前所有问题**，并且是**后续继续写代码的唯一基准**。

```text
/internal
├── domain                        # 领域层（不关心怎么跑）
│   ├── common
│   │   └── enums.go              # 极少量跨域枚举（Side / Direction）
│   │
│   ├── market
│   │   ├── event.go              # MarketEvent（行情事件）
│   │   └── market.go             # Bar / Tick / Snapshot
│   │
│   ├── instrument
│   │   └── instrument.go         # 股票 / 期货 / 期权抽象
│   │
│   ├── strategy
│   │   ├── strategy.go           # Strategy 接口
│   │   └── signal.go             # Strategy Signal / Intent
│   │
│   ├── order
│   │   └── order.go              # 订单领域模型
│   │
│   ├── trade
│   │   └── trade.go              # 成交领域模型
│   │
│   └── portfolio
│       ├── portfolio.go          # 投资组合
│       └── position.go           # 持仓
│
├── application                   # 应用层（怎么把 domain 跑起来）
│   ├── engine
│   │   ├── strategy
│   │   │   ├── engine.go         # Strategy Engine 门面
│   │   │   ├── dispatcher.go     # 策略调度（并发 / 隔离）
│   │   │   ├── runtime.go        # 策略运行时
│   │   │   └── registry.go       # 策略注册
│   │   │
│   │   ├── risk                  # Risk Engine（后续）
│   │   ├── execution             # Execution Engine（后续）
│   │   └── backtest              # Backtest Engine（后续）
│   │
│   └── service
│       └── trading_app.go        # 系统装配与编排
│
├── infrastructure                # 外部世界
│   ├── marketdata                # 行情接入
│   ├── broker                    # 券商 / 模拟撮合
│   ├── strategy                  # 具体策略实现
│   │   └── examples
│   │       └── ma_cross.go
│   └── storage                   # 存储
│
└── interfaces                    # 对外接口
    ├── grpc
    └── cli
```

> ⚠️ **以后任何文件，都只能落在这棵树上**
> ⚠️ 如果放不进去，说明设计还没想清楚

---

# 三、Strategy Engine 的【最终权威定位】

## 1️⃣ Strategy（Domain）

```go
// internal/domain/strategy/strategy.go

type Strategy interface {
    Name() string
    OnInit(ctx Context) error
    OnMarketEvent(ctx Context, event market.Event) ([]Signal, error)
    OnStop(ctx Context) error
}
```

它只表达一件事：

> **“一个策略，在收到行情时，可能会产生交易意图”**

---

## 2️⃣ Strategy Engine（Application / Engine）

**唯一入口文件：**

```text
/internal/application/engine/strategy/engine.go
```

### 它做什么

* 生命周期管理（Start / Stop）
* 行情入口（OnMarketEvent）
* 把事件交给 Dispatcher

### 它绝对不做什么

* 不算指标
* 不做下单
* 不做风控
* 不关心策略内容

### 它的本质

> **一个“守门员 + 门面（Facade）”**

---

## 3️⃣ Dispatcher / Runtime 的角色（澄清之前的混乱）

| 组件         | 职责        |
| ---------- | --------- |
| Engine     | 对外接口 + 边界 |
| Dispatcher | 并发调度 / 隔离 |
| Runtime    | 单策略运行上下文  |

之前看起来“重复”，是因为我在**不同抽象层反复解释它们**，现在这一次统一收敛。

---
