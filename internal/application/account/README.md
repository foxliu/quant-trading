# Account Context

## 职责边界（冻结）

Account Context：
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

Account Snapshot 是唯一读接口
