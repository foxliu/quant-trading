-- 002_seed_runtime_profile_example.sql
-- 用途：写入/更新一个可运行的示例 runtime profile（手动执行）
-- 数据库：SQLite

INSERT INTO t_trader_runtime_profile (
    profile_key,
    active,
    broker_name,
    account_id,
    initial_cash,
    ctp_front_addr,
    ctp_broker_id,
    ctp_user_id,
    ctp_investor_id,
    ctp_password,
    strategies_json
) VALUES (
    'default',
    1,
    'CTP',
    'paper_account_001',
    1000000,
    'tcp://180.168.146.187:10101',
    '9999',
    'your_user_id',
    'your_investor_id',
    'your_password',
    '[{"name":"MA_Cross_IH2503","type":"ma_cross","params":{"short_period":5,"long_period":20}},{"name":"Breakout_IH2503","type":"breakout","params":{"period":20,"quantity":100}}]'
)
ON CONFLICT(profile_key) DO UPDATE SET
    active = excluded.active,
    broker_name = excluded.broker_name,
    account_id = excluded.account_id,
    initial_cash = excluded.initial_cash,
    ctp_front_addr = excluded.ctp_front_addr,
    ctp_broker_id = excluded.ctp_broker_id,
    ctp_user_id = excluded.ctp_user_id,
    ctp_investor_id = excluded.ctp_investor_id,
    ctp_password = excluded.ctp_password,
    strategies_json = excluded.strategies_json;

-- 可选：确保只有一个 active 配置（若有多条，请按需调整）
-- UPDATE t_trader_runtime_profile
-- SET active = CASE WHEN profile_key = 'default' THEN 1 ELSE 0 END;
