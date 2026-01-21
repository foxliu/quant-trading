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

--- 
## 整体架构中的位置
```text
Market Event
   ↓
Strategy Engine
   ↓
Signal
   ↓
Planner
   ↓
Order            ← 业务订单（意图清晰）
   ↓
Risk Engine
   ↓
Execution Engine ← 当前部分
   ↓
Broker / Exchange / Simulator

```
---
## 执行流完整时序
```text
Strategy
   │
   │ SubmitOrderIntent
   ▼
StrategyContext
   │
   ▼
AccountContext
   │   ├─ 风控 / 账户选择
   │   └─ build ExecutionRequest
   ▼
ExecutionContext
   │   ├─ submit to broker
   │   ├─ manage order state
   │   └─ emit ExecutionEvent
   ▼
Position / Account Update

```
---
## Execution 与 Position / Risk 的最终关系图
```text
Risk Rule
   ↓
Risk Engine
   ↓ Result
Risk Coordinator
   ↓ Command
Execution Controller
   ├── Cancel Orders
   ├── Read Position
   └── Market Close
         ↓
Execution Event
         ↓
Position Context
         ↓
PnL / Account / Risk
```
---
## Execution Order 状态机
```text
NEW
 ↓ submit
SUBMITTED
 ↓ accepted
ACCEPTED
 ↓ fill
PARTIALLY_FILLED ──┐
 ↓ fill            │
FILLED             │
                   │
        cancel ───▶ CANCELED
        reject ───▶ REJECTED
```