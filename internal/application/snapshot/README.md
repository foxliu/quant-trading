### 1️⃣ Snapshot 不是 Event
* Snapshot 不是业务驱动
* Snapshot 不参与决策
* Snapshot 不进入 Event Bus

### 2️⃣ Snapshot 是“系统加速器”

* 只用于：
    * Replay 加速
    * Debug 快速定位
    * 回测切片

### 3️⃣ Snapshot 必须是：
* 明确所属 Context
* 明确时间点
* 明确可序列化