## 「PnL Context + PnL Engine（事件驱动，非侵入）」

> 目标：
> **在不破坏 Position Context 的前提下，构建一个可被 Risk / Monitor / UI 消费的 PnL 计算链路**

---

# 一、PnL 设计原则（冻结）

### 1️⃣ PnL 不“驱动”任何行为

* 不下单
* 不调仓
* 不修正 Position
  👉 **纯观察者（Observer）**

---

### 2️⃣ PnL 分两类，必须严格区分

| 类型                 | 定义   | 触发源                |
| ------------------ | ---- | ------------------ |
| **Realized PnL**   | 已经落袋 | Execution Event    |
| **Unrealized PnL** | 浮动盈亏 | Market Price Event |

二者**不能混算、不能推导替代**

---

### 3️⃣ PnL Engine 不直接依赖 Strategy / Planner

依赖关系冻结为：

```
Execution Engine ──▶ Position Context
                      │
                      ▼
                 PnL Context
                      │
            Market Price Context
```

---

# 一、先给冻结级结论（非常重要）

> **PnL 重估不是 Context**
> **PnL 重估是一个“只读合成引擎（Projection Engine）”**

它的职责是：

* **不维护状态**
* **不接收 Event**
* **不写回任何 Context**

它只是把 **Account + Position + Market Price**
在某一个时间点 **合成为一个确定性的视图**

---

# 二、为什么 PnL Engine 不能是 Context

这是你之后一定会被问到的问题，我现在一次性冻结答案。

### 如果 PnL 是 Context，会发生什么？

| 问题          | 后果           |
| ----------- | ------------ |
| 依赖 Position | 状态重复         |
| 依赖 Account  | 双写风险         |
| 接收 Event    | Replay 顺序地狱  |
| 有 Snapshot  | Snapshot 不一致 |

**PnL 是“派生值”，不是“源状态”**

---

# 三、PnL Revaluation 的冻结输入 / 输出

## 输入（只读）

```text
- Account Balance
- Position（可为 nil）
- Market Price（最新 / 指定时间）
```

## 输出（一次性）

```text
- UnrealizedPnL
- Equity
- Exposure
```

---

# 四、冻结工程位置

```text
internal/
 └─ application/
    └─ pnl/
       ├─ engine.go        // 主引擎（冻结）
       └─ snapshot.go     // 可选快照（用于回放加速）
```

> ❗注意
> **它不在 domain，不在 context**

---

# 五、完整冻结合成路径（非常重要）

```text
Execution Event
   ↓
Position Context        → Qty / AvgPrice
Account Context         → Cash / Balance
Market Price Context    → Latest Price
   ↓
PnL Engine (Revaluate)
   ↓
Equity / PnL / Exposure
```

> **PnL Engine 永远是“最后一步”**

---

# 六、Replay / Backtest 的一致性（你已经赢了）

因为你现在：

* Position 可 Replay
* Account 可 Replay
* Market Price 可 Replay
* PnL 无状态

你自动获得了：

> **Live == Backtest == Replay 的数值一致性**

这是**绝大多数交易系统失败的地方**，而你已经绕过去了。

---

# 七、冻结规则总结（请记住）

1. **PnL 是派生值，不是状态**
2. **派生值不进入 Event Bus**
3. **派生值不做 Snapshot（除非为了加速）**

---

