## 调用链
```text
Market Event / Bar
        ↓
Strategy Engine
        ↓
Strategy (OnBar / OnTick)
        ↓
Signal
        ↓
Planner   ← 这里需要注意
        ↓
Order(s)
        ↓
Risk Engine
        ↓
Execution / Broker Adapter

```
---
## Risk Coordinator + Force Close 执行链路

**把 Risk Rule Engine 的“判罚结果”，可靠、可控、可扩展地转化为真实的执行动作，** 并且：
* 不污染 Risk Engine
* 不侵入 Strategy / Planner
* 不绕过 Execution Engine
* 可以长期冻结

```text
┌────────────┐
│ Risk Rules │
└─────┬──────┘
      │ Result
┌─────▼──────────┐
│ Risk Engine    │
└─────┬──────────┘
      │ Result Channel
┌─────▼──────────────┐
│ Risk Coordinator   │ 
│  - 状态机          │
│  - 幂等控制        │
│  - 动作翻译        │
└─────┬──────────────┘
      │ Command
┌─────▼──────────────┐
│ Execution Engine   │
│  - Close Position  │
│  - Cancel Orders   │
└────────────────────┘
```
**👉 Risk Engine 只“说话”，Coordinator 才“动手”。**

**Risk Coordinator 只做 4 件事：**
1. 监听 RiskResult
2. 做 去重 / 幂等
3. 将 Action 翻译为 Execution Command
4. 串行触发执行（避免连环强平） 
**它不判断风险、不计算仓位、不看策略**