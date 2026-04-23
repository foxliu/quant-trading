package snapshot

import "time"

type Snapshot interface {
	Name() string // Position / RiskContext / Account ...
	Timestamp() time.Time
}
