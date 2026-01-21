package execution

type OrderManager interface {
	// CancelAll 取消订单
	CancelAll(symbol string) error
}
