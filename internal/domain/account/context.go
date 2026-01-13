package account

type Context interface {
	// GetAccount 账户信息
	GetAccount() Account

	// GetAvailableCash 现金
	GetAvailableCash() float64

	// GetPosition 位置
	GetPosition(symbol string) (Position, bool)
	GetAllPositions() []Position

	// PlaceOrder 订单抽象
	PlaceOrder(symbol string, side OrderSide, price float64, qty int64) (Order, error)
}
