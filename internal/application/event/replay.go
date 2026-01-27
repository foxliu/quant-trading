package event

import "quant-trading/internal/application/snapshot"

/*
Replay 的核心思想

Replay 不是“重放函数”，而是“重新驱动整个系统”
*/

type ReplayEngine struct {
	bus   Bus
	store snapshot.Store
}

func NewReplayer(bus Bus, store snapshot.Store) *ReplayEngine {
	return &ReplayEngine{bus: bus, store: store}
}

func (r *ReplayEngine) ReplayFromSnapshot(sn snapshot.Snapshot, events []*Envelope) {
	if sn != nil {
		restoreAll(sn)
	}

	for _, evt := range events {
		if evt.Timestamp.After(snapshot.Timestamp()) {
			r.bus.Publish(evt)
		}
	}
}
