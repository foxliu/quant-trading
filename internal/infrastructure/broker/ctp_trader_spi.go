package broker

import (
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/execution"
	"quant-trading/internal/infrastructure/logger"
	"time"

	"github.com/pseudocodes/go2ctp/thost"
	"go.uber.org/zap"
)

var log = logger.Logger.With(zap.String("module", "broker.ctp_trader_spi"))

// 内部SPI回调处理器（go2ctp风格）
type ctpTraderSpi struct {
	adapter *CTPAdapter
}

func (s *ctpTraderSpi) OnRsqUserLogin(
	req *thost.CThostFtdcReqUserLoginField,
	info *thost.CThostFtdcRspInfoField,
	reqID int,
	isLast bool,
) {
	if info.ErrorID != 0 {
		log.Error("CTP 登录失败", zap.String("ErrorMsg", thost.BytesToString(info.ErrorMsg[:])))
		return
	}
	tradingDay := thost.BytesToString(req.TradingDay[:])
	s.adapter.SetTradingDay(tradingDay)
}

func (s *ctpTraderSpi) OnFrontConnected() {
	// 登录
	log.Info("CTP 前置已连接，开始登录...")
	req := thost.CThostFtdcReqUserLoginField{
		BrokerID: s.adapter.brokerID,
		UserID:   s.adapter.userID,
		Password: s.adapter.password,
	}
	s.adapter.traderApi.ReqUserLogin(&req, 0)
}

// OnRtnOrder / OnRtnTrade 等回调推送到 events 通道（省略部分，生产可完整实现）

func (s *ctpTraderSpi) OnFrontDisconnected(nReason int) {
	log.Info("CTP 前置断开连接", zap.Int("nReason", nReason))
	s.adapter.events <- execution.Event{
		Type:      execution.EventDisconnected,
		Timestamp: time.Now(),
	}
}

func (s *ctpTraderSpi) OnHeartBeatWarning(nTimeLapse int) {
	log.Info("CTP 心跳警告", zap.Int("nTimelapse", nTimeLapse))
}

func (s *ctpTraderSpi) OnRspAuthenticate(pRspAuthenticateField *thost.CThostFtdcRspAuthenticateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspUserLogin(pRspUserLogin *thost.CThostFtdcRspUserLoginField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	if pRspInfo.ErrorID != 0 {
		log.Error("CTP 登录失败", zap.Int32("ErrorID", int32(pRspInfo.ErrorID)), zap.ByteString("ErrorMsg", pRspInfo.ErrorMsg[:]))
		return
	}
	log.Info("✅ CTP 登录成功")
	// 自动查询持仓和资金
	s.adapter.traderApi.ReqQryInvestorPosition(nil, 0)
	s.adapter.traderApi.ReqQryTradingAccount(nil, 0)
}

func (s *ctpTraderSpi) OnRspUserLogout(pUserLogout *thost.CThostFtdcUserLogoutField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspUserPasswordUpdate(pUserPasswordUpdate *thost.CThostFtdcUserPasswordUpdateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspTradingAccountPasswordUpdate(pTradingAccountPasswordUpdate *thost.CThostFtdcTradingAccountPasswordUpdateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspUserAuthMethod(pRspUserAuthMethod *thost.CThostFtdcRspUserAuthMethodField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspGenUserCaptcha(pRspGenUserCaptcha *thost.CThostFtdcRspGenUserCaptchaField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspGenUserText(pRspGenUserText *thost.CThostFtdcRspGenUserTextField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspOrderInsert(pInputOrder *thost.CThostFtdcInputOrderField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	if pRspInfo.ErrorID != 0 {
		log.Warn("CTP 订单提交失败", zap.Int32("ErrorID", int32(pRspInfo.ErrorID)), zap.ByteString("ErrorMsg", pRspInfo.ErrorMsg[:]))
		s.adapter.events <- execution.Event{
			Type:      execution.EventOrderRejected,
			OrderID:   thost.BytesToString(pInputOrder.OrderRef[:]),
			Timestamp: time.Now(),
			Reason:    thost.BytesToString(pRspInfo.ErrorMsg[:]),
		}
		return
	}
}

func (s *ctpTraderSpi) OnRspParkedOrderInsert(pParkedOrder *thost.CThostFtdcParkedOrderField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspParkedOrderAction(pParkedOrderAction *thost.CThostFtdcParkedOrderActionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspOrderAction(pInputOrderAction *thost.CThostFtdcInputOrderActionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryMaxOrderVolume(pQryMaxOrderVolume *thost.CThostFtdcQryMaxOrderVolumeField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspSettlementInfoConfirm(pSettlementInfoConfirm *thost.CThostFtdcSettlementInfoConfirmField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspRemoveParkedOrder(pRemoveParkedOrder *thost.CThostFtdcRemoveParkedOrderField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspRemoveParkedOrderAction(pRemoveParkedOrderAction *thost.CThostFtdcRemoveParkedOrderActionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspExecOrderInsert(pInputExecOrder *thost.CThostFtdcInputExecOrderField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspExecOrderAction(pInputExecOrderAction *thost.CThostFtdcInputExecOrderActionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspForQuoteInsert(pInputForQuote *thost.CThostFtdcInputForQuoteField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQuoteInsert(pInputQuote *thost.CThostFtdcInputQuoteField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQuoteAction(pInputQuoteAction *thost.CThostFtdcInputQuoteActionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspBatchOrderAction(pInputBatchOrderAction *thost.CThostFtdcInputBatchOrderActionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspOptionSelfCloseInsert(pInputOptionSelfClose *thost.CThostFtdcInputOptionSelfCloseField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspOptionSelfCloseAction(pInputOptionSelfCloseAction *thost.CThostFtdcInputOptionSelfCloseActionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspCombActionInsert(pInputCombAction *thost.CThostFtdcInputCombActionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryOrder(pOrder *thost.CThostFtdcOrderField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryTrade(pTrade *thost.CThostFtdcTradeField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorPosition(pInvestorPosition *thost.CThostFtdcInvestorPositionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	if pRspInfo.ErrorID != 0 || pInvestorPosition == nil {
		log.Error("CTP 获取持仓失败: %d %s\n", zap.Int32("ErrorID", int32(pRspInfo.ErrorID)), zap.ByteString("ErrorMsg", pRspInfo.ErrorMsg[:]))
		return
	}
	//symbol := pInvestorPosition.InstrumentID.String()
	//s.adapter.mu.Lock()
	//s.adapter.positions[symbol] = trade.Position{
	//	Instrument: instrument.Instrument{Symbol: symbol},
	//	Symbol:     symbol,
	//	Qty:        int64(pInvestorPosition.Position),
	//}
	//s.adapter.mu.Lock()
}

func (s *ctpTraderSpi) OnRspQryTradingAccount(pTradingAccount *thost.CThostFtdcTradingAccountField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	if pRspInfo != nil && pRspInfo.ErrorID != 0 {
		log.Error("CTP 获取账户失败: %d %s\n", zap.Int32("ErrorID", int32(pRspInfo.ErrorID)), zap.ByteString("ErrorMsg", pRspInfo.ErrorMsg[:]))
		return
	}
	log.Info("CTP 获取账户成功", zap.Any("Account", pTradingAccount))

	evtLoop := &event.Envelope{
		Payload:   pTradingAccount,
		Source:    "CTP",
		Type:      event.EventCTPTradingAccountRtn,
		Timestamp: time.Now(),
	}
	s.adapter.eventBus.Publish(evtLoop)
}

func (s *ctpTraderSpi) OnRspQryInvestor(pInvestor *thost.CThostFtdcInvestorField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryTradingCode(pTradingCode *thost.CThostFtdcTradingCodeField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInstrumentMarginRate(pInstrumentMarginRate *thost.CThostFtdcInstrumentMarginRateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInstrumentCommissionRate(pInstrumentCommissionRate *thost.CThostFtdcInstrumentCommissionRateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryUserSession(pUserSession *thost.CThostFtdcUserSessionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryExchange(pExchange *thost.CThostFtdcExchangeField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryProduct(pProduct *thost.CThostFtdcProductField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInstrument(pInstrument *thost.CThostFtdcInstrumentField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryDepthMarketData(pDepthMarketData *thost.CThostFtdcDepthMarketDataField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryTraderOffer(pTraderOffer *thost.CThostFtdcTraderOfferField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySettlementInfo(pSettlementInfo *thost.CThostFtdcSettlementInfoField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryTransferBank(pTransferBank *thost.CThostFtdcTransferBankField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorPositionDetail(pInvestorPositionDetail *thost.CThostFtdcInvestorPositionDetailField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryNotice(pNotice *thost.CThostFtdcNoticeField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySettlementInfoConfirm(pSettlementInfoConfirm *thost.CThostFtdcSettlementInfoConfirmField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorPositionCombineDetail(pInvestorPositionCombineDetail *thost.CThostFtdcInvestorPositionCombineDetailField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryCFMMCTradingAccountKey(pCFMMCTradingAccountKey *thost.CThostFtdcCFMMCTradingAccountKeyField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryEWarrantOffset(pEWarrantOffset *thost.CThostFtdcEWarrantOffsetField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorProductGroupMargin(pInvestorProductGroupMargin *thost.CThostFtdcInvestorProductGroupMarginField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryExchangeMarginRate(pExchangeMarginRate *thost.CThostFtdcExchangeMarginRateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryExchangeMarginRateAdjust(pExchangeMarginRateAdjust *thost.CThostFtdcExchangeMarginRateAdjustField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryExchangeRate(pExchangeRate *thost.CThostFtdcExchangeRateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySecAgentACIDMap(pSecAgentACIDMap *thost.CThostFtdcSecAgentACIDMapField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryProductExchRate(pProductExchRate *thost.CThostFtdcProductExchRateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryProductGroup(pProductGroup *thost.CThostFtdcProductGroupField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryMMInstrumentCommissionRate(pMMInstrumentCommissionRate *thost.CThostFtdcMMInstrumentCommissionRateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryMMOptionInstrCommRate(pMMOptionInstrCommRate *thost.CThostFtdcMMOptionInstrCommRateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInstrumentOrderCommRate(pInstrumentOrderCommRate *thost.CThostFtdcInstrumentOrderCommRateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySecAgentTradingAccount(pTradingAccount *thost.CThostFtdcTradingAccountField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySecAgentCheckMode(pSecAgentCheckMode *thost.CThostFtdcSecAgentCheckModeField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySecAgentTradeInfo(pSecAgentTradeInfo *thost.CThostFtdcSecAgentTradeInfoField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryOptionInstrTradeCost(pOptionInstrTradeCost *thost.CThostFtdcOptionInstrTradeCostField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryOptionInstrCommRate(pOptionInstrCommRate *thost.CThostFtdcOptionInstrCommRateField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryExecOrder(pExecOrder *thost.CThostFtdcExecOrderField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryForQuote(pForQuote *thost.CThostFtdcForQuoteField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryQuote(pQuote *thost.CThostFtdcQuoteField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryOptionSelfClose(pOptionSelfClose *thost.CThostFtdcOptionSelfCloseField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestUnit(pInvestUnit *thost.CThostFtdcInvestUnitField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryCombInstrumentGuard(pCombInstrumentGuard *thost.CThostFtdcCombInstrumentGuardField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryCombAction(pCombAction *thost.CThostFtdcCombActionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryTransferSerial(pTransferSerial *thost.CThostFtdcTransferSerialField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryAccountregister(pAccountregister *thost.CThostFtdcAccountregisterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspError(pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Error("❌ CTP通用错误", zap.Int32("ErrorID", int32(pRspInfo.ErrorID)), zap.ByteString("ErrorMsg", pRspInfo.ErrorMsg[:]))
}

func (s *ctpTraderSpi) OnRtnOrder(pOrder *thost.CThostFtdcOrderField) {
	s.adapter.eventBus.Publish(&event.Envelope{
		Type:      event.EventCTPOrderRtn,
		Source:    "CTP",
		Timestamp: time.Now(),
		Payload:   pOrder,
	})
}

func (s *ctpTraderSpi) OnRtnTrade(pTrade *thost.CThostFtdcTradeField) {
	log.Info("完成交易", zap.ByteString("成交回报", pTrade.InstrumentID[:]),
		zap.Float64("价格", float64(pTrade.Price)),
		zap.Float64("数量", float64(pTrade.Volume)),
	)

	// 推送成交事件
	//s.adapter.events <- execution.Event{
	//	Type:      execution.EventOrderFilled,
	//	OrderID:   string(pTrade.OrderRef[:]),
	//	Price:     float64(pTrade.Price),
	//	FilledQty: int64(pTrade.Volume),
	//	UpdateTime: time.Now(),
	//}
	// 触发 ApplyFill（通过 AccountContext，由 execution engine 处理）
	// 这里只推送事件，实际 ApplyFill 在 paperEngine 或 dispatcher 中监听事件后调用
}

func (s *ctpTraderSpi) OnErrRtnOrderInsert(pInputOrder *thost.CThostFtdcInputOrderField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnOrderAction(pOrderAction *thost.CThostFtdcOrderActionField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnInstrumentStatus(pInstrumentStatus *thost.CThostFtdcInstrumentStatusField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnBulletin(pBulletin *thost.CThostFtdcBulletinField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnTradingNotice(pTradingNoticeInfo *thost.CThostFtdcTradingNoticeInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnErrorConditionalOrder(pErrorConditionalOrder *thost.CThostFtdcErrorConditionalOrderField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnExecOrder(pExecOrder *thost.CThostFtdcExecOrderField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnExecOrderInsert(pInputExecOrder *thost.CThostFtdcInputExecOrderField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnExecOrderAction(pExecOrderAction *thost.CThostFtdcExecOrderActionField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnForQuoteInsert(pInputForQuote *thost.CThostFtdcInputForQuoteField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnQuote(pQuote *thost.CThostFtdcQuoteField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnQuoteInsert(pInputQuote *thost.CThostFtdcInputQuoteField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnQuoteAction(pQuoteAction *thost.CThostFtdcQuoteActionField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnForQuoteRsp(pForQuoteRsp *thost.CThostFtdcForQuoteRspField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnCFMMCTradingAccountToken(pCFMMCTradingAccountToken *thost.CThostFtdcCFMMCTradingAccountTokenField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnBatchOrderAction(pBatchOrderAction *thost.CThostFtdcBatchOrderActionField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnOptionSelfClose(pOptionSelfClose *thost.CThostFtdcOptionSelfCloseField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnOptionSelfCloseInsert(pInputOptionSelfClose *thost.CThostFtdcInputOptionSelfCloseField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnOptionSelfCloseAction(pOptionSelfCloseAction *thost.CThostFtdcOptionSelfCloseActionField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnCombAction(pCombAction *thost.CThostFtdcCombActionField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnCombActionInsert(pInputCombAction *thost.CThostFtdcInputCombActionField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryContractBank(pContractBank *thost.CThostFtdcContractBankField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryParkedOrder(pParkedOrder *thost.CThostFtdcParkedOrderField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryParkedOrderAction(pParkedOrderAction *thost.CThostFtdcParkedOrderActionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryTradingNotice(pTradingNotice *thost.CThostFtdcTradingNoticeField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryBrokerTradingParams(pBrokerTradingParams *thost.CThostFtdcBrokerTradingParamsField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryBrokerTradingAlgos(pBrokerTradingAlgos *thost.CThostFtdcBrokerTradingAlgosField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQueryCFMMCTradingAccountToken(pQueryCFMMCTradingAccountToken *thost.CThostFtdcQueryCFMMCTradingAccountTokenField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnFromBankToFutureByBank(pRspTransfer *thost.CThostFtdcRspTransferField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnFromFutureToBankByBank(pRspTransfer *thost.CThostFtdcRspTransferField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnRepealFromBankToFutureByBank(pRspRepeal *thost.CThostFtdcRspRepealField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnRepealFromFutureToBankByBank(pRspRepeal *thost.CThostFtdcRspRepealField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnFromBankToFutureByFuture(pRspTransfer *thost.CThostFtdcRspTransferField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnFromFutureToBankByFuture(pRspTransfer *thost.CThostFtdcRspTransferField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnRepealFromBankToFutureByFutureManual(pRspRepeal *thost.CThostFtdcRspRepealField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnRepealFromFutureToBankByFutureManual(pRspRepeal *thost.CThostFtdcRspRepealField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnQueryBankBalanceByFuture(pNotifyQueryAccount *thost.CThostFtdcNotifyQueryAccountField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnBankToFutureByFuture(pReqTransfer *thost.CThostFtdcReqTransferField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnFutureToBankByFuture(pReqTransfer *thost.CThostFtdcReqTransferField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnRepealBankToFutureByFutureManual(pReqRepeal *thost.CThostFtdcReqRepealField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnRepealFutureToBankByFutureManual(pReqRepeal *thost.CThostFtdcReqRepealField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnQueryBankBalanceByFuture(pReqQueryAccount *thost.CThostFtdcReqQueryAccountField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnRepealFromBankToFutureByFuture(pRspRepeal *thost.CThostFtdcRspRepealField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnRepealFromFutureToBankByFuture(pRspRepeal *thost.CThostFtdcRspRepealField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspFromBankToFutureByFuture(pReqTransfer *thost.CThostFtdcReqTransferField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspFromFutureToBankByFuture(pReqTransfer *thost.CThostFtdcReqTransferField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQueryBankAccountMoneyByFuture(pReqQueryAccount *thost.CThostFtdcReqQueryAccountField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnOpenAccountByBank(pOpenAccount *thost.CThostFtdcOpenAccountField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnCancelAccountByBank(pCancelAccount *thost.CThostFtdcCancelAccountField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnChangeAccountByBank(pChangeAccount *thost.CThostFtdcChangeAccountField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryClassifiedInstrument(pInstrument *thost.CThostFtdcInstrumentField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryCombPromotionParam(pCombPromotionParam *thost.CThostFtdcCombPromotionParamField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRiskSettleInvstPosition(pRiskSettleInvstPosition *thost.CThostFtdcRiskSettleInvstPositionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRiskSettleProductStatus(pRiskSettleProductStatus *thost.CThostFtdcRiskSettleProductStatusField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySPBMFutureParameter(pSPBMFutureParameter *thost.CThostFtdcSPBMFutureParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySPBMOptionParameter(pSPBMOptionParameter *thost.CThostFtdcSPBMOptionParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySPBMIntraParameter(pSPBMIntraParameter *thost.CThostFtdcSPBMIntraParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySPBMInterParameter(pSPBMInterParameter *thost.CThostFtdcSPBMInterParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySPBMPortfDefinition(pSPBMPortfDefinition *thost.CThostFtdcSPBMPortfDefinitionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySPBMInvestorPortfDef(pSPBMInvestorPortfDef *thost.CThostFtdcSPBMInvestorPortfDefField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorPortfMarginRatio(pInvestorPortfMarginRatio *thost.CThostFtdcInvestorPortfMarginRatioField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorProdSPBMDetail(pInvestorProdSPBMDetail *thost.CThostFtdcInvestorProdSPBMDetailField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorCommoditySPMMMargin(pInvestorCommoditySPMMMargin *thost.CThostFtdcInvestorCommoditySPMMMarginField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorCommodityGroupSPMMMargin(pInvestorCommodityGroupSPMMMargin *thost.CThostFtdcInvestorCommodityGroupSPMMMarginField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySPMMInstParam(pSPMMInstParam *thost.CThostFtdcSPMMInstParamField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySPMMProductParam(pSPMMProductParam *thost.CThostFtdcSPMMProductParamField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQrySPBMAddOnInterParameter(pSPBMAddOnInterParameter *thost.CThostFtdcSPBMAddOnInterParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRCAMSCombProductInfo(pRCAMSCombProductInfo *thost.CThostFtdcRCAMSCombProductInfoField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRCAMSInstrParameter(pRCAMSInstrParameter *thost.CThostFtdcRCAMSInstrParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRCAMSIntraParameter(pRCAMSIntraParameter *thost.CThostFtdcRCAMSIntraParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRCAMSInterParameter(pRCAMSInterParameter *thost.CThostFtdcRCAMSInterParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRCAMSShortOptAdjustParam(pRCAMSShortOptAdjustParam *thost.CThostFtdcRCAMSShortOptAdjustParamField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRCAMSInvestorCombPosition(pRCAMSInvestorCombPosition *thost.CThostFtdcRCAMSInvestorCombPositionField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorProdRCAMSMargin(pInvestorProdRCAMSMargin *thost.CThostFtdcInvestorProdRCAMSMarginField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRULEInstrParameter(pRULEInstrParameter *thost.CThostFtdcRULEInstrParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRULEIntraParameter(pRULEIntraParameter *thost.CThostFtdcRULEIntraParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryRULEInterParameter(pRULEInterParameter *thost.CThostFtdcRULEInterParameterField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorProdRULEMargin(pInvestorProdRULEMargin *thost.CThostFtdcInvestorProdRULEMarginField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorPortfSetting(pInvestorPortfSetting *thost.CThostFtdcInvestorPortfSettingField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryInvestorInfoCommRec(pInvestorInfoCommRec *thost.CThostFtdcInvestorInfoCommRecField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryCombLeg(pCombLeg *thost.CThostFtdcCombLegField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspOffsetSetting(pInputOffsetSetting *thost.CThostFtdcInputOffsetSettingField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspCancelOffsetSetting(pInputOffsetSetting *thost.CThostFtdcInputOffsetSettingField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRtnOffsetSetting(pOffsetSetting *thost.CThostFtdcOffsetSettingField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnOffsetSetting(pInputOffsetSetting *thost.CThostFtdcInputOffsetSettingField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnErrRtnCancelOffsetSetting(pCancelOffsetSetting *thost.CThostFtdcCancelOffsetSettingField, pRspInfo *thost.CThostFtdcRspInfoField) {
	//TODO implement me
	panic("implement me")
}

func (s *ctpTraderSpi) OnRspQryOffsetSetting(pOffsetSetting *thost.CThostFtdcOffsetSettingField, pRspInfo *thost.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//TODO implement me
	panic("implement me")
}
