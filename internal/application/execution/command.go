package execution

import "time"

type CommandType string

const (
	CommandForceClose CommandType = "FORCE_CLOSE"
	CommandCancelAll  CommandType = "CANCEL_ALL"
)

type Command struct {
	Type   CommandType
	Symbol string

	Reason string
	Time   time.Time
}
