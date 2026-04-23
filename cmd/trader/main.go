package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	aAccount "quant-trading/internal/application/account"
	"quant-trading/internal/application/capital"
	execpaper "quant-trading/internal/application/execution/paper"
	"quant-trading/internal/application/event"
	"quant-trading/internal/application/portfolio"
	"quant-trading/internal/application/risk"
	"quant-trading/internal/application/runtime"
	dStrategy "quant-trading/internal/domain/strategy"
	"quant-trading/internal/infrastructure/broker"
	"quant-trading/internal/infrastructure/config"
	"quant-trading/internal/infrastructure/db"
	"quant-trading/internal/infrastructure/logger"
	"quant-trading/internal/infrastructure/persistence"
	"quant-trading/internal/infrastructure/strategy/examples"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	configPath := flag.String("config", "configs/trader.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.LoadTraderConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "配置加载失败: %v\n", err)
		os.Exit(1)
	}
	if err := logger.InitLogger(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	log := zap.L().With(zap.String("module", "cmd.trader"))
	log.Info("=== 量化交易系统 - CTP 实盘交易模式（集成 RiskEngine） ===", zap.String("config", *configPath))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := cfg.DB.DSN
	if dsn == "" {
		dsn = "trader.db"
	}
	sqliteDB := db.InitSQLite(dsn)

	accountRepo := persistence.NewRepository(sqliteDB)
	bus := event.NewMemoryBus()

	accountSvc := aAccount.NewService(accountRepo, bus)
	accountManager := aAccount.NewManager(accountSvc)

	persister := aAccount.NewPersister(accountRepo, bus)
	persister.Start()

	profile, err := persistence.GetActiveTraderRuntimeProfile(sqliteDB)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatal("加载运行时配置失败", zap.Error(err))
		}
		log.Warn("未找到数据库运行时配置，尝试从配置文件迁移一次")
		if err := bootstrapRuntimeProfileFromYAML(sqliteDB, cfg); err != nil {
			log.Fatal("运行时配置缺失，且迁移失败；请先在数据库写入 t_trader_runtime_profile", zap.Error(err))
		}
		profile, err = persistence.GetActiveTraderRuntimeProfile(sqliteDB)
		if err != nil {
			log.Fatal("迁移后加载运行时配置失败", zap.Error(err))
		}
	}

	initialCash := profile.InitialCash
	if initialCash <= 0 {
		initialCash = 100_000
	}

	capEngine := capital.NewMemoryEngine(initialCash)
	portEngine := portfolio.NewMemoryEngine()

	stableKey := profile.AccountID
	if stableKey == "" {
		log.Fatal("数据库运行时配置 account_id 不能为空")
	}

	acc, err := accountSvc.LoadOrCreateAccount(aAccount.LoadAccountCommand{
		UserID:     "user-001",
		BrokerName: profile.BrokerName,
		Alias:      stableKey,
		Config: map[string]any{
			"initial_capital": initialCash,
		},
	})
	if err != nil {
		log.Fatal("加载或创建账户失败", zap.Error(err))
	}

	accountCtx, err := accountManager.GetContext(acc.ID, capEngine, portEngine)
	if err != nil {
		log.Fatal("获取账户上下文失败", zap.Error(err))
	}

	refresh := runtime.NewRefreshScheduler(accountCtx, bus)
	refresh.Start()

	log.Info("账户初始化完成",
		zap.String("account_id", acc.ID.String()),
		zap.Float64("available", accountCtx.Available()),
		zap.Float64("equity", accountCtx.Equity()))

	strategies, err := profile.Strategies()
	if err != nil {
		log.Fatal("解析数据库策略配置失败", zap.Error(err))
	}
	if len(strategies) == 0 {
		log.Fatal("数据库策略配置为空，请在 t_trader_runtime_profile.strategies_json 中配置")
	}

	fmt.Println("创建 CTP Broker...")
	ctpBroker, err := broker.NewCTPAdapter(
		profile.CTPFrontAddr,
		profile.CTPBrokerID,
		profile.CTPInvestorID,
		profile.CTPUserID,
		profile.CTPPassword,
		profile.AccountID,
		accountCtx,
		portEngine,
	)
	if err != nil {
		log.Fatal("CTP 连接失败", zap.Error(err))
	}
	log.Info("CTP Broker 已登录成功")

	executionEngine := execpaper.NewEngine(ctpBroker, accountCtx)
	log.Info("执行引擎已就绪")

	riskEngine := risk.NewEngine()
	log.Info("风控引擎已初始化")

	scheduler := runtime.NewScheduler(accountCtx, riskEngine, bus, executionEngine)

	for _, spec := range strategies {
		stg, err := buildStrategy(spec)
		if err != nil {
			log.Fatal("构建策略失败", zap.String("strategy", spec.Name), zap.Error(err))
		}
		if err := scheduler.RegisterStrategy(stg); err != nil {
			log.Fatal("策略注册失败", zap.String("strategy", spec.Name), zap.Error(err))
		}
	}

	if err := scheduler.Start(ctx); err != nil {
		log.Fatal("Scheduler 启动失败", zap.Error(err))
	}

	coordinator := risk.NewCoordinator(riskEngine, bus)
	if err := coordinator.Start(ctx); err != nil {
		log.Warn("风控协调器启动失败", zap.Error(err))
	}

	log.Info("✅ CTP 实盘系统已启动！RiskEngine 已集成",
		zap.String("config_file", *configPath),
		zap.String("ctp_account_id", profile.AccountID),
		zap.Float64("initial_cash", initialCash))

	<-ctx.Done()
	scheduler.Stop()
	log.Info("交易系统已停止")
}

func bootstrapRuntimeProfileFromYAML(sqliteDB *gorm.DB, cfg *config.TraderConfig) error {
	if cfg.CTP.AccountID == "" || cfg.CTP.FrontAddr == "" {
		return fmt.Errorf("缺少用于迁移的 ctp 关键信息")
	}
	strategies := make([]persistence.StrategySpec, 0, len(cfg.Strategies))
	for _, s := range cfg.Strategies {
		strategies = append(strategies, persistence.StrategySpec{
			Name:   s.Name,
			Type:   s.Type,
			Params: s.Params,
		})
	}
	if len(strategies) == 0 {
		strategies = append(strategies, persistence.StrategySpec{
			Name: "MA_Cross_IH2503",
			Type: "ma_cross",
			Params: map[string]any{
				"short_period": 5,
				"long_period":  20,
			},
		})
	}
	payload, err := json.Marshal(strategies)
	if err != nil {
		return err
	}
	initialCash := cfg.Account.InitialCash
	if initialCash <= 0 {
		initialCash = 100_000
	}
	return persistence.UpsertTraderRuntimeProfile(sqliteDB, &persistence.TraderRuntimeProfile{
		Key:            "default",
		Active:         true,
		BrokerName:     "CTP",
		AccountID:      cfg.CTP.AccountID,
		InitialCash:    initialCash,
		CTPFrontAddr:   cfg.CTP.FrontAddr,
		CTPBrokerID:    cfg.CTP.BrokerID,
		CTPUserID:      cfg.CTP.UserID,
		CTPInvestorID:  cfg.CTP.InvestorID,
		CTPPassword:    cfg.CTP.Password,
		StrategiesJSON: string(payload),
	})
}

func buildStrategy(spec persistence.StrategySpec) (dStrategy.Strategy, error) {
	name := strings.TrimSpace(spec.Name)
	if name == "" {
		return nil, fmt.Errorf("strategy name 不能为空")
	}

	switch strings.ToLower(strings.TrimSpace(spec.Type)) {
	case "ma_cross", "macross":
		shortPeriod := paramInt(spec.Params, "short_period", 5)
		longPeriod := paramInt(spec.Params, "long_period", 20)
		return examples.NewMACrossStrategy(name, shortPeriod, longPeriod), nil
	case "breakout":
		period := paramInt(spec.Params, "period", 20)
		quantity := paramFloat(spec.Params, "quantity", 100)
		return examples.NewBreakoutStrategy(name, period, quantity), nil
	default:
		return nil, fmt.Errorf("不支持的策略类型: %s", spec.Type)
	}
}

func paramInt(params map[string]interface{}, key string, def int) int {
	if params == nil {
		return def
	}
	v, ok := params[key]
	if !ok {
		return def
	}
	switch n := v.(type) {
	case int:
		return n
	case int8:
		return int(n)
	case int16:
		return int(n)
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float32:
		return int(n)
	case float64:
		return int(n)
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return int(i)
		}
	case string:
		if i, err := strconv.Atoi(n); err == nil {
			return i
		}
	}
	return def
}

func paramFloat(params map[string]interface{}, key string, def float64) float64 {
	if params == nil {
		return def
	}
	v, ok := params[key]
	if !ok {
		return def
	}
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int8:
		return float64(n)
	case int16:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case json.Number:
		if f, err := n.Float64(); err == nil {
			return f
		}
	case string:
		if f, err := strconv.ParseFloat(n, 64); err == nil {
			return f
		}
	}
	return def
}
