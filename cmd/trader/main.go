package main

import (
	"context"
	"flag"
	"fmt"
	aAccount "quant-trading/internal/application/account"
	"quant-trading/internal/application/capital"
	"quant-trading/internal/application/event"
	"quant-trading/internal/application/paper"
	"quant-trading/internal/application/portfolio"
	"quant-trading/internal/application/risk" // 新增
	"quant-trading/internal/application/runtime"
	"quant-trading/internal/domain/account"
	"quant-trading/internal/infrastructure/broker"
	"quant-trading/internal/infrastructure/config"
	"quant-trading/internal/infrastructure/logger"
	"quant-trading/internal/infrastructure/strategy/examples"

	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config", "configs/trader.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.LoadTraderConfig(*configPath)
	if err != nil {
		return
	}
	defer logger.Sync()

	log := logger.Logger.With(zap.String("module", "cmd.trader"))
	log.Info("=== 量化交易系统 - CTP 实盘交易模式（集成 RiskEngine） ===", zap.String("config", *configPath))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 加载配置
	fmt.Println("1. 加载配置...")
	cfg, err = config.LoadTraderConfig(*configPath)
	if err != nil {
		log.Fatal("配置加载失败", zap.Error(err))
	}

	// 2. 创建账户组件
	fmt.Println("2. 初始化账户组件...")
	domainAcc := &account.Account{AccountID: cfg.CTP.AccountID}
	capEngine := capital.NewMemoryEngine(cfg.Account.InitialCash)
	portEngine := portfolio.NewMemoryEngine()

	accountCtx := aAccount.NewContext(domainAcc, capEngine, portEngine)
	log.Info("账户初始化完成",
		zap.Float64("available", accountCtx.Available()),
		zap.Float64("equity", accountCtx.Equity()))

	// 3. 创建 CTP Broker
	fmt.Println("3. 创建 CTP Broker...")
	ctpBroker, err := broker.NewCTPAdapter(
		cfg.CTP.FrontAddr,
		cfg.CTP.BrokerID,
		cfg.CTP.UserID,
		cfg.CTP.InvestorID,
		cfg.CTP.Password,
		cfg.CTP.AccountID,
	)
	if err != nil {
		log.Fatal("CTP 连接失败", zap.Error(err))
	}
	log.Info("CTP Broker 已登录成功")

	// 4. 创建执行引擎
	fmt.Println("4. 创建执行引擎...")
	executionEngine := paper.NewEngine(ctpBroker, accountCtx)
	log.Info("执行引擎已就绪")

	// 5. 创建风控引擎（新增）
	fmt.Println("5. 创建 RiskEngine...")
	riskEngine := risk.NewEngine()
	log.Info("风控引擎已初始化")

	// 6. 创建事件总线
	bus := event.NewMemoryBus()

	// 7. 创建多策略调度器（注入 RiskEngine）
	fmt.Println("6. 创建 Scheduler...")
	scheduler := runtime.NewScheduler(accountCtx, riskEngine, bus, executionEngine)

	// 8. 注册策略
	fmt.Println("7. 注册策略...")
	maStrategy := examples.NewMACrossStrategy("MA_Cross_IH2503", 5, 20)
	if err := scheduler.RegisterStrategy(maStrategy); err != nil {
		log.Fatal("策略注册失败", zap.Error(err))
	}

	// 9. 启动 Scheduler
	fmt.Println("8. 启动多策略调度器...")
	if err := scheduler.Start(ctx); err != nil {
		log.Fatal("Scheduler 启动失败", zap.Error(err))
	}

	// 10. 启动风控协调器（可选，监听事件）
	coordinator := risk.NewCoordinator(riskEngine, bus)
	if err := coordinator.Start(ctx); err != nil {
		log.Warn("风控协调器启动失败", zap.Error(err))
	}

	log.Info("✅ CTP 实盘系统已启动！RiskEngine 已集成",
		zap.String("config_file", *configPath),
		zap.String("account_id", cfg.CTP.AccountID),
		zap.Float64("initial_cash", cfg.Account.InitialCash))

	// 保持运行
	<-ctx.Done()
	scheduler.Stop()
	log.Info("交易系统已停止")
}
