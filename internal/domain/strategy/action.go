package strategy

type Action string

const (
	ActionBuy   Action = "BUY"
	ActionSell  Action = "SELL"
	ActionClose Action = "ClOSE"
	ActionHold  Action = "HOLD"
)
