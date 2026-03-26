package execution

/*
EventType
=========

Execution Event 是“执行事实”，不是业务意图
*/
type EventType int

const (
	EventOrderSubmitted EventType = iota
	EventOrderFilled
	EventOrderCanceled
	EventOrderRejected
	EventDisconnected
	EventOrderAccepted
	EventOrderPartiallyFilled
)

func (t EventType) String() string {
	switch t {
	case EventOrderFilled:
		return "FILLED"
	case EventOrderCanceled:
		return "CANCELED"
	case EventOrderRejected:
		return "REJECTED"
	case EventDisconnected:
		return "DISCONNECTED"
	case EventOrderAccepted:
		return "ACCEPTED"
	case EventOrderPartiallyFilled:
		return "PARTIALLY_FILLED"
	case EventOrderSubmitted:
		return "SUBMITTED"
	default:
		return "UNKNOWN"
	}
}
