# Execution Engine

## 职责边界（冻结）

Execution Engine 只负责：
- 将 Order 提交到执行系统
- 产生 Execution Event

Execution Engine 不负责：
- Signal 解析
- 风控
- 仓位计算
- 业务决策
- 同步成交保证

## 数据流
```text
Order
↓
Execution Engine
↓
Execution Event
↓
Position / Account / Order State
```

## 设计原则

- Execution 是事实源，不是状态源
- Execution Event 是系统唯一的成交事实
- 所有仓位 / PnL / Exposure 必须由 Event 驱动
