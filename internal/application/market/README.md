# Market Price Context（最终冻结版本）

> 目标：提供**多 Symbol**、**可 Snapshot**、**事件驱动 + 拉取并存**的统一市场价格上下文，作为 Strategy / Execution 等 Context 的只读依赖。

---

## 一、设计边界与职责（冻结）

**职责**

* 维护多 Symbol 的最新行情（Quote / Bar / Trade 的最小必要子集）
* 提供**一致性 Snapshot**（跨 Symbol、跨字段）
* 接收 Market Event（Tick / Bar / Quote）并更新内部状态
* 提供只读查询接口（同步拉取）

**不做的事**

* 不负责行情订阅/网络 IO
* 不做指标计算（指标属于 Indicator Context）
* 不对外暴露可变结构（严格只读）

---

## 二、核心概念模型

### 1. Symbol

* 全局唯一键（如：`BINANCE:BTCUSDT:SPOT`）

### 2. PriceState（单 Symbol）

最小但完整的价格状态集合：

* Last Price
* Bid / Ask
* Volume（累计）
* OHLC（最近 Bar）
* EventTime（行情时间）
* UpdateSeq（单调递增，用于一致性）

### 3. MarketSnapshot（跨 Symbol）

* 不可变
* 捕获某一逻辑时间点的**全量 Symbol 视图**
* 用于回测 / 决策一致性

---

## 三、对外接口（冻结）

```go
// MarketPriceContext 只读接口
// Strategy / Execution 只能依赖该接口

type MarketPriceContext interface {
    // 单 Symbol 查询
    Get(symbol string) (PriceView, bool)

    // 多 Symbol 查询
    GetMany(symbols []string) map[string]PriceView

    // 当前一致性快照（全量）
    Snapshot() MarketSnapshot

    // 快照（子集）
    SnapshotOf(symbols []string) MarketSnapshot
}
```

---

## 四、只读视图对象（冻结）

```go
// PriceView 为不可变值对象
// Strategy 永远拿不到指针

type PriceView struct {
    Symbol     string
    Last       float64
    Bid        float64
    Ask        float64
    Volume     float64

    Open       float64
    High       float64
    Low        float64
    Close      float64

    EventTime  int64
    UpdateSeq  uint64
}
```

---

## 五、Snapshot 定义（冻结）

```go
// MarketSnapshot 是跨 Symbol 的一致性视图
// 一旦创建，不可修改

type MarketSnapshot struct {
    Seq        uint64            // Snapshot 序列号
    AsOf      int64             // 逻辑时间
    Prices    map[string]PriceView
}
```

**一致性语义**

* 同一个 Snapshot 内：所有 PriceView 均 <= Snapshot.Seq
* Snapshot 为深拷贝，Strategy 可安全长期持有

---

## 六、内部可变状态（实现细节，冻结结构）

```go
// priceState 为内部可变结构，不对外暴露

type priceState struct {
    symbol    string

    last      float64
    bid       float64
    ask       float64
    volume    float64

    open      float64
    high      float64
    low       float64
    close     float64

    eventTime int64
    updateSeq uint64
}
```

---

## 七、事件驱动更新模型（冻结）

```go
// MarketEvent 由 Event Bus 投递

type MarketEvent interface {
    Symbol() string
    EventTime() int64
}

// 示例：TickEvent / BarEvent / QuoteEvent
```

**更新规则**

* 单 Symbol 内严格单调递增 updateSeq
* 乱序 EventTime 允许，但 updateSeq 永远递增
* 最新事件覆盖对应字段（Bar 覆盖 OHLC，Tick 覆盖 Last）

---

## 八、并发与一致性策略（冻结）

* 内部使用 `RWMutex`
* 写路径：事件驱动（Event Bus → ApplyEvent）
* 读路径：

    * `Get` / `GetMany`：RLock + 值拷贝
    * `Snapshot`：RLock + 全量深拷贝

**保证**

* Snapshot 获取期间无写入
* Snapshot 与 Get 不互相阻塞（读读并发）

---

## 九、最小实现骨架（冻结）

```go
// marketPriceContext 为唯一实现

type marketPriceContext struct {
    mu      sync.RWMutex
    seq     uint64
    prices  map[string]*priceState
}

func NewMarketPriceContext() *marketPriceContext {
    return &marketPriceContext{
        prices: make(map[string]*priceState),
    }
}
```

---

## 十、ApplyEvent（冻结语义）

```go
func (m *marketPriceContext) ApplyEvent(evt MarketEvent) {
    m.mu.Lock()
    defer m.mu.Unlock()

    ps, ok := m.prices[evt.Symbol()]
    if !ok {
        ps = &priceState{symbol: evt.Symbol()}
        m.prices[evt.Symbol()] = ps
    }

    m.seq++
    ps.updateSeq = m.seq
    ps.eventTime = evt.EventTime()

    // 根据事件类型更新字段（略）
}
```

---

## 十一、Snapshot 实现语义（冻结）

```go
func (m *marketPriceContext) Snapshot() MarketSnapshot {
    m.mu.RLock()
    defer m.mu.RUnlock()

    snap := MarketSnapshot{
        Seq:   m.seq,
        AsOf: time.Now().UnixNano(),
        Prices: make(map[string]PriceView, len(m.prices)),
    }

    for sym, ps := range m.prices {
        snap.Prices[sym] = toView(ps)
    }
    return snap
}
```

---

## 十二、与其他 Context 的边界（冻结）

* **Strategy Context**：

    * 只依赖 `MarketPriceContext` 接口
    * 每次决策使用 Snapshot

* **Execution Context**：

    * 可读取实时 PriceView（非 Snapshot）

* **Indicator Context**：

    * 订阅 Market Event，自行维护时间序列

---

## 十三、冻结结论

该版本：

* 支持 **多 Symbol**
* 支持 **一致性 Snapshot**
* 明确 **只读边界**
* 可安全用于：回测 / 实盘 / 并发策略

> 后续演进只能新增：字段、事件类型、性能优化；
> **禁止破坏接口与一致性语义。**
