package broker

import (
	"context"
	"fmt"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/trade"
	"sync"
	"time"

	"github.com/pseudocodes/go2ctp/ctp"
	"github.com/pseudocodes/go2ctp/thost"
)

// CTPAdapter 上期CTP真实适配器（实现 Broker 接口）
type CTPAdapter struct {
	mu        sync.RWMutex
	traderApi thost.TraderApi
	events    chan execution.Event

	// 缓存（供GetPositions / GetBalance使用）
	positions map[string]trade.Position // symbol -> Position
	balance   float64
	equity    float64

	// 订单映射（CTP OrderRef -> 本地 Order）
	orderMap map[string]*order.Order

	// 配置
	accountID  string
	brokerID   thost.TThostFtdcBrokerIDType
	userID     thost.TThostFtdcUserIDType
	investorID thost.TThostFtdcInvestorIDType
	password   thost.TThostFtdcPasswordType
	frontAddr  string // 交易前置地址
	symbol     string // 默认合约（可动态）
}

// NewCTPAdapter 创建CTP适配器
// frontAddr 示例："tcp://180.168.146.187:10000" （测试环境）或券商生产地址
func NewCTPAdapter(frontAddr, brokerID, investorID, userID, password, accountID string) (*CTPAdapter, error) {
	var inID thost.TThostFtdcInvestorIDType
	copy(inID[:], investorID)

	var bID thost.TThostFtdcBrokerIDType
	copy(bID[:], brokerID)

	var pwd thost.TThostFtdcPasswordType
	copy(pwd[:], password)

	var uID thost.TThostFtdcUserIDType
	copy(uID[:], userID)

	a := &CTPAdapter{
		events:    make(chan execution.Event, 200),
		positions: make(map[string]trade.Position),
		orderMap:  make(map[string]*order.Order),

		accountID:  accountID,
		investorID: inID,
		userID:     uID,
		password:   pwd,
		brokerID:   bID,
		frontAddr:  frontAddr,
		symbol:     "IH2503", // 默认合约，可后续从 Order 动态设置
	}

	// 创建 TraderApi
	a.traderApi = ctp.CreateTraderApi(ctp.TraderFlowPath("./ctp_flow/")) //流文件目录（自动创建）

	// 注册 SPI （回调）
	spi := &ctpTraderSpi{adapter: a}
	a.traderApi.RegisterSpi(spi)

	// 连接前置
	a.traderApi.RegisterFront(a.frontAddr)
	a.traderApi.Init()

	// 等待登录（实际生产建议加超时/重连）
	time.Sleep(2 * time.Second) // 简化，生产用WaitGroup或channel

	return a, nil
}

// SubmitOrder 实现Broker 适配器接口
func (a *CTPAdapter) SubmitOrder(ctx context.Context, ord *order.Order) (string, error) {
	a.mu.Lock()
	a.orderMap[ord.ID()] = ord
	a.mu.Unlock()

	var instrumentID thost.TThostFtdcInstrumentIDType
	copy(instrumentID[:], ord.Symbol())

	var orderRef thost.TThostFtdcOrderRefType
	copy(orderRef[:], ord.ID())
	req := thost.CThostFtdcInputOrderField{
		BrokerID:            a.brokerID,
		InvestorID:          a.investorID,
		InstrumentID:        instrumentID,
		OrderRef:            orderRef,
		OrderPriceType:      thost.THOST_FTDC_OPT_LimitPrice, // 支持市价/限价转换
		Direction:           convertDirection(ord.Side()),
		CombOffsetFlag:      thost.TThostFtdcCombOffsetFlagType{'0'}, // 开仓
		LimitPrice:          thost.TThostFtdcPriceType(ord.Price()),
		VolumeTotalOriginal: thost.TThostFtdcVolumeType(ord.Qty()),
	}

	errCode := a.traderApi.ReqOrderInsert(&req, 0)
	if errCode != 0 {
		return "", fmt.Errorf("CTP 下单失败，code: %d", errCode)
	}

	a.events <- execution.Event{Type: execution.EventOrderSubmitted, OrderID: ord.ID(), Timestamp: time.Now()}
	return ord.ID(), nil
}

// CancelOrder 实现Broker 适配器接口
func (a *CTPAdapter) CancelOrder(ctx context.Context, orderID string) error {
	// 简化实现（实际需保存 OrderRef + FrontID + SessionID）
	var oRef thost.TThostFtdcOrderRefType
	copy(oRef[:], orderID)
	req := thost.CThostFtdcInputOrderActionField{
		BrokerID:   a.brokerID,
		InvestorID: a.investorID,
		OrderRef:   oRef,
		ActionFlag: thost.THOST_FTDC_AF_Delete,
	}
	errCode := a.traderApi.ReqOrderAction(&req, 0)
	if errCode != 0 {
		return fmt.Errorf("CTP 撤单失败，code: %d", errCode)
	}
	return nil
}

// GetPositions / GetBalance（通过 ReqQryInvestorPosition / ReqQryTradingAccount 实现，回调处理）
func (a *CTPAdapter) GetPositions(ctx context.Context) ([]trade.Position, error) {
	// 实际通过回调 OnRspQryInvestorPosition 填充，简化返回空（生产需缓存）
	a.mu.RLock()
	defer a.mu.RUnlock()
	pos := make([]trade.Position, 0, len(a.positions))
	for _, p := range a.positions {
		pos = append(pos, p)
	}
	return pos, nil
}
func (a *CTPAdapter) GetBalance(ctx context.Context) (cash float64, equity float64, err error) {
	// 通过回调 OnRspQryTradingAccount 返回
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.balance, a.equity, nil
}

func (a *CTPAdapter) SubscribeEvents(ctx context.Context) <-chan execution.Event {
	return a.events
}

func convertDirection(side order.Side) thost.TThostFtdcDirectionType {
	if side == order.Buy {
		return thost.THOST_FTDC_D_Buy
	}
	return thost.THOST_FTDC_D_Sell
}
