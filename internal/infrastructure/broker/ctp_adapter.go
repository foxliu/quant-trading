package broker

import (
	"context"
	"fmt"
	"quant-trading/internal/application/account"
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/domain/order"
	"quant-trading/internal/domain/portfolio"
	"sync"
	"time"

	"github.com/pseudocodes/go2ctp/ctp"
	"github.com/pseudocodes/go2ctp/thost"
)

// CTPAdapter 上期CTP真实适配器（实现 Broker 接口）
type CTPAdapter struct {
	mu         sync.RWMutex
	traderApi  thost.TraderApi
	events     chan execution.Event
	eventBus   event.Bus
	accountCtx *account.Context
	portfolio  portfolio.Engine
	// 配置
	accountID string
	brokerID  thost.TThostFtdcBrokerIDType // 经纪商代码
	userID    thost.TThostFtdcUserIDType   // 用户登录名
	// 投资者代码, 就是在期货公司开立的资金账号（交易账号），是资金结算和持仓归属的唯一标识
	investorID thost.TThostFtdcInvestorIDType
	password   thost.TThostFtdcPasswordType
	// frontAddr 交易前置地址， 期货公司提供的 CTP 服务器网络地址， CTP 分为行情前置（MdFront）和交易前置（TdFront），
	// 两者地址不同。登录时需要分别连接对应的地址。
	frontAddr  string
	symbol     string // 默认合约（可动态）
	tradingDay string

	// 穿透式监管
	appID    thost.TThostFtdcAppIDType    // 客户端应用标识
	authCode thost.TThostFtdcAuthCodeType // 授权码/认证码

	pending  map[CTPReqID]*PendingRequest
	reqIDGen CTPReqID
}

// NewCTPAdapter 创建CTP适配器
// frontAddr 示例："tcp://180.168.146.187:10000" （测试环境）或券商生产地址
func NewCTPAdapter(
	frontAddr, brokerID, investorID, userID, password, accountID string,
	accountCtx *account.Context,
	portfolio portfolio.Engine,
) (*CTPAdapter, error) {
	var inID thost.TThostFtdcInvestorIDType
	copy(inID[:], investorID)

	var bID thost.TThostFtdcBrokerIDType
	copy(bID[:], brokerID)

	var pwd thost.TThostFtdcPasswordType
	copy(pwd[:], password)

	var uID thost.TThostFtdcUserIDType
	copy(uID[:], userID)

	a := &CTPAdapter{
		events:     make(chan execution.Event, 200),
		accountCtx: accountCtx,
		portfolio:  portfolio,
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
func (a *CTPAdapter) CancelOrder(ctx context.Context, ord *order.Order) error {
	// 简化实现（实际需保存 OrderRef + FrontID + SessionID）
	var oRef thost.TThostFtdcOrderRefType
	copy(oRef[:], ord.ID())
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

func (a *CTPAdapter) QueryOrderStatus(ctx context.Context, ord *order.Order) (*order.Order, error) {
	var oRef thost.TThostFtdcOrderSysIDType
	copy(oRef[:], ord.ID())
	req := thost.CThostFtdcQryOrderField{
		BrokerID:   a.brokerID,
		InvestorID: a.investorID,
		OrderSysID: oRef,
	}
}

// GetPositions / QueryAccount（通过 ReqQryInvestorPosition / ReqQryTradingAccount 实现，回调处理）
func (a *CTPAdapter) GetPositions(ctx context.Context) ([]portfolio.Position, error) {
	// 实际通过回调 OnRspQryInvestorPosition 填充，简化返回空（生产需缓存）
	a.mu.RLock()
	defer a.mu.RUnlock()
	positions, _ := a.accountCtx.GetPositions()
	pos := make([]portfolio.Position, 0, len(positions))
	for _, p := range positions {
		pos = append(pos, p)
	}
	return pos, nil
}
func (a *CTPAdapter) QueryAccount(ctx context.Context) (*thost.CThostFtdcTradingAccountField, error) {
	reqID := a.nextReqID()
	pr := &PendingRequest{
		ch: make(chan *execution.AccountEvent, 1),
	}
	// 通过回调 OnRspQryTradingAccount 返回
	a.mu.Lock()
	a.pending[reqID] = pr
	a.mu.Unlock()
	req := &thost.CThostFtdcQryTradingAccountField{
		BrokerID:   a.brokerID,
		InvestorID: a.investorID,
	}
	errCode := a.traderApi.ReqQryTradingAccount(req, reqID.Value())
	if errCode != 0 {
		return nil, fmt.Errorf("CTP 查询账户失败，code: %d", errCode)
	}
	// 等待完成事件
	select {
	case evt := <-pr.ch:
		a.cleanupPending(reqID)
		return evt.Data, evt.Err
	case <-ctx.Done():
		a.cleanupPending(reqID)
		return nil, ctx.Err()
	}
}

func (a *CTPAdapter) cleanupPending(reqID CTPReqID) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.pending, reqID)
}

func (a *CTPAdapter) SubscribeEvents(ctx context.Context) <-chan execution.Event {
	return a.events
}

func (a *CTPAdapter) SetTradingDay(tradingDay string) {
	a.mu.Lock()
	a.tradingDay = tradingDay
	a.mu.Unlock()
}

func (a *CTPAdapter) StartEventLoop() {
	a.eventBus.Subscribe(event.EventCTPOrderRtn, a.handleAccountQueryEvent)
}

func (a *CTPAdapter) nextReqID() CTPReqID {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.reqIDGen++
	return a.reqIDGen
}

func (a *CTPAdapter) handleAccountQueryEvent(e *event.Envelope) {
	switch e.Type {
	case event.EventCTPOrderRtn:
		if evt, ok := e.Payload.(*execution.AccountEvent); ok {
			a.mu.Lock()
			reqID := CTPReqID(evt.ReqID)
			pr, ok := a.pending[reqID]
			if !ok {
				return
			}
			if evt.IsLast {
				pr.ch <- evt
			}
		}
	default:
		return
	}
}

func convertDirection(side order.Side) thost.TThostFtdcDirectionType {
	if side == order.Buy {
		return thost.THOST_FTDC_D_Buy
	}
	return thost.THOST_FTDC_D_Sell
}
