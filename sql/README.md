# SQL 手动脚本目录

本目录用于存放**所有需要人工执行**的 SQL 脚本。

## 约定

- 所有手动 SQL 必须放在 `sql/` 下，禁止散落到其他目录。
- 文件命名采用有序前缀：`NNN_描述.sql`，例如 `001_init_xxx.sql`。
- 每个 SQL 文件头部必须说明：用途、适用数据库、是否可重复执行。
- 优先写成幂等语句（如 `IF NOT EXISTS`、`ON CONFLICT`）。
- 变更上线前，先在测试库验证，再在生产库执行。

## 当前脚本

- `001_init_trader_runtime_profile.sql`：初始化运行时配置表。
- `002_seed_runtime_profile_example.sql`：插入/更新示例运行时配置。
