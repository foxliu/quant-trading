# Position Context

## 核心原则（冻结）

- Position 只由 Execution Event 推导
- 不读取 Signal / Order / Planner
- 是系统中唯一的仓位事实来源

## 状态机规则

- 同方向成交 → 加仓
- 反方向成交 → 减仓 / 平仓 / 反手
- Qty = 0 → Position 不存在

## 下游依赖

- Account Context
- Risk Exposure
- PnL 计算
