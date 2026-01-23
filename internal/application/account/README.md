# Account Context

## 职责边界（冻结）

Account Context：
- 基于 Execution Event，确定性地维护资金状态
- 聚合 Position Snapshot
- 维护账户资金状态
- 计算账户级 Equity / Exposure

不做：
- 下单决策
- 风控判断
- 成交解析

## 数据来源

Execution Event  → Position Context  
Position Snapshot → Account Context

## 对外出口

- Account Snapshot 是唯一读接口
- 只输出状态（Snapshot / Query）
  - Cash Balance
  - Equity
  - Realized PnL
  - Unrealized PnL（来自 Position + Market）

---
## 事实链路
```text
Strategy
   ↓ Signal
Planner
   ↓ Order
Risk Engine
   ↓ Approved Order
Execution Engine
   ↓ Execution Event
Position Context
   ↓ Position Snapshot
Account Context
```
---
## 资金链路
```text
Market Event
   ↓
Strategy → Signal
   ↓
Risk / Planner
   ↓
Order
   ↓
Execution Engine
   ↓
Execution Event
   ↓
Position Context   → 仓位
Account Context    → 资金
Market Price       → 重估
   ↓
Equity / PnL
```