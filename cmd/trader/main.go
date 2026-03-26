package main

import (
	"context"
	"fmt"
	"os"
	"quant-trading/internal/application/event"
	"quant-trading/internal/application/runtime"
	"quant-trading/internal/infrastructure/logger"
	"quant-trading/internal/infrastructure/strategy/examples"

	"quant-trading/internal/application/account"
	"quant-trading/internal/application/capital"
	"quant-trading/internal/application/paper"
	"quant-trading/internal/application/portfolio"
	dAccount "quant-trading/internal/domain/account" // domain Account
	"quant-trading/internal/infrastructure/broker"

	"go.uber.org/zap"
)

func main() {
	err := logger.InitLogger()
	if err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
	defer logger.Sync()

	log := logger.Logger.With(zap.String("module", "cmd.trader"))

	log.Info("=== 量化交易系统 - CTP 实盘交易模式（上期所期货） ===")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 从环境变量读取 CTP 配置（安全）
	frontAddr := os.Getenv("CTP_FRONT_ADDR")
	brokerID := os.Getenv("CTP_BROKER_ID")
	investorID := os.Getenv("CTP_INVESTOR_ID")
	userID := os.Getenv("CTP_USER_ID")
	password := os.Getenv("CTP_PASSWORD")
	accountID := os.Getenv("CTP_ACCOUNT_ID")

	if frontAddr == "" || brokerID == "" || investorID == "" || password == "" {
		log.Fatal("❌ 请设置 CTP 环境变量: CTP_FRONT_ADDR / CTP_BROKER_ID / CTP_INVESTOR_ID / CTP_PASSWORD")
	}

	// 2. 创建领域 Account + capital + portfolio（必传参数）
	log.Info("1. 初始化账户组件...")
	domainAcc := &dAccount.Account{AccountID: accountID}
	capEngine := capital.NewMemoryEngine(1000000.0) // 期货保证金示例 100万
	portEngine := portfolio.NewMemoryEngine()

	accountCtx := account.NewContext(domainAcc, capEngine, portEngine)
	log.Info("完成初始化账户组件", zap.Float64("初始可用资金", accountCtx.Available()), zap.Float64("权益", accountCtx.Equity()))

	// 3. 创建 CTP 真实 Broker（替换 PaperBroker）
	log.Info("2. 创建 CTP 适配器（pseudocodes/go2ctp）...")
	ctpBroker, err := broker.NewCTPAdapter(frontAddr, brokerID, investorID, userID, password, accountID)
	if err != nil {
		log.Fatal("❌ CTP 连接失败: %v", zap.Error(err))
	}
	log.Info("   CTP Broker 已登录成功（实时交易就绪）")

	// 4. 创建执行引擎（复用 paper 包，未来可改名为 execution）
	log.Info("3. 创建执行引擎...")
	executionEngine := paper.NewEngine(ctpBroker, accountCtx)
	log.Info("执行引擎已就绪（CTP 模式）")

	// 5. 创建事件总线
	bus := event.NewMemoryBus()

	// 6. 加载策略（支持多策略扩展）
	log.Info("4. 创建 Scheduler...")
	scheduler := runtime.NewScheduler(accountCtx, bus, executionEngine)

	// 7. 注册策略
	log.Info("5. 注册策略...")
	maStrategy := examples.NewMACrossStrategy("MA_Cross_IH2503", 5, 20)
	if err := scheduler.RegisterStrategy(maStrategy); err != nil {
		log.Fatal("注册策略失败: %v", zap.Error(err))
	}

	// 8. 启动系统
	log.Info("6. 启动 CTP 实盘交易系统...")
	if err := executionEngine.Start(ctx); err != nil {
		log.Fatal("启动失败: %v", zap.Error(err))
	}

	log.Info("✅ 系统已启动！CTP 实盘模式运行中...")
	log.Info("当前合约示例: IH2503（可动态切换）")
	log.Info("按 Ctrl+C 停止交易")

	// 保持运行（实际生产中接 runtime/dispatcher + 事件总线）
	<-ctx.Done()
	log.Info("交易系统已停止")
}
