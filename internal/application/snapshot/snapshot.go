package snapshot

import "time"

type Snapshot interface {
	Name() string // Position / Risk / Account ...
	Timestamp() time.Time
}
