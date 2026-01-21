package risk

type Action string

const (
	ActionNone        Action = "NONE"
	ActionRejectOrder Action = "REJECT_ORDER"
	ActionForceClose  Action = "FORCE_CLOSE"
	ActionHaltTrading Action = "HALT_TRADING"
)
