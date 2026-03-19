package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"quant-trading/internal/infrastructure/logger"

	"quant-trading/internal/application/account"
	"quant-trading/internal/application/capital"
	"quant-trading/internal/application/paper"
	"quant-trading/internal/application/portfolio"
	dAccount "quant-trading/internal/domain/account" // domain Account
	"quant-trading/internal/infrastructure/broker"
)

func main() {

	err := logger.InitLogger()
	if err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	defer logger.Sync()

	fmt.Println("=== 量化交易系统 - CTP 实盘交易模式（上期所期货） ===")

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
	fmt.Println("1. 初始化账户组件...")
	domainAcc := &dAccount.Account{AccountID: accountID}
	capEngine := capital.NewMemoryEngine(1000000.0) // 期货保证金示例 100万
	portEngine := portfolio.NewMemoryEngine()

	accountCtx := account.NewContext(domainAcc, capEngine, portEngine)
	fmt.Printf("   初始可用资金: %.2f | 权益: %.2f\n\n", accountCtx.Available(), accountCtx.Equity())

	// 3. 创建 CTP 真实 Broker（替换 PaperBroker）
	fmt.Println("2. 创建 CTP 适配器（pseudocodes/go2ctp）...")
	ctpBroker, err := broker.NewCTPAdapter(frontAddr, brokerID, investorID, userID, password, accountID)
	if err != nil {
		log.Fatalf("❌ CTP 连接失败: %v", err)
	}
	fmt.Print("   CTP Broker 已登录成功（实时交易就绪）\n\n")

	// 4. 创建执行引擎（复用 paper 包，未来可改名为 execution）
	fmt.Println("3. 创建执行引擎...")
	executionEngine := paper.NewEngine(ctpBroker, accountCtx)
	fmt.Print("   执行引擎已就绪（CTP 模式）\n\n")

	// 5. 加载策略（支持多策略扩展）
	fmt.Println("4. 加载双均线策略...")
	//maStrategy := examples.NewMACrossStrategy("MA_Cross_CTP_IH2503", 5, 20)

	// 6. 启动系统
	fmt.Println("5. 启动 CTP 实盘交易系统...")
	if err := executionEngine.Start(ctx); err != nil {
		log.Fatalf("启动失败: %v", err)
	}

	fmt.Println("   ✅ 系统已启动！CTP 实盘模式运行中...")
	fmt.Println("   当前合约示例: IH2503（可动态切换）")
	fmt.Println("   按 Ctrl+C 停止交易")

	// 保持运行（实际生产中接 runtime/dispatcher + 事件总线）
	<-ctx.Done()
	fmt.Println("交易系统已停止")
}
