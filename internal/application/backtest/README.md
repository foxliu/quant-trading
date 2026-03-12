## Backtest Engine 的真实执行拓扑
```text
BacktestEngine
      │
      ▼
RuntimeClock
      │
      ▼
MarketEvent
      │
      ▼
Strategy.OnBar()
      │
      ▼
ExecutionService
      │
      ▼
Simulator
      │
      ▼
FillEvent
      │
      ▼
AccountContext
      │
      ▼
Portfolio / Capital
```
Clock的职责
```text
推进时间
产生 MarketEvent
驱动策略执行
```