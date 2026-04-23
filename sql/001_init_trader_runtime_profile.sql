-- 001_init_trader_runtime_profile.sql
-- 用途：初始化运行时配置表（手动执行）
-- 说明：仅在需要手工建表或核对结构时执行。

CREATE TABLE IF NOT EXISTS t_trader_runtime_profile (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_key      VARCHAR(64)  NOT NULL UNIQUE,
    active           BOOLEAN      NOT NULL DEFAULT 1,
    broker_name      VARCHAR(32)  NOT NULL DEFAULT 'CTP',
    account_id       VARCHAR(64)  NOT NULL,
    initial_cash     REAL         NOT NULL DEFAULT 100000,
    ctp_front_addr   VARCHAR(255) NOT NULL,
    ctp_broker_id    VARCHAR(64)  NOT NULL,
    ctp_user_id      VARCHAR(64)  NOT NULL,
    ctp_investor_id  VARCHAR(64)  NOT NULL,
    ctp_password     VARCHAR(255) NOT NULL,
    strategies_json  TEXT         NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_trader_runtime_profile_active
    ON t_trader_runtime_profile(active);
