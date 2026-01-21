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
