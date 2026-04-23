
# Position Context 的工程定位（冻结）

## 1️⃣ Position Context 在整体架构中的位置

```text
Strategy
  └── Signal (TargetQty)
        ↓
Planner
  └── Order (未修正，Quantity = Target 占位)
        ↓
Position Engine   ←【我们现在做的】
  ├── Position Context（仓位事实）
  ├── Target → Delta 计算
  └── 修正 Order.Quantity / Side
        ↓
Risk / Execution（后续）
```

---

## 2️⃣ Position Context 的职责（冻结）

**它只负责三件事：**

1. 维护“当前仓位事实”
2. 根据目标仓位计算变化量（Delta）
3. 决定 **最终买 / 卖方向**

**它不负责：**

* 风控（风险限额、杠杆）
* 下单
* 成交回报
* 资金校验

---
