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
	// TODO: 实现快照恢复逻辑
	// if sn != nil {
	//     restoreAll(sn)
	// }

	for _, evt := range events {
		// TODO: 实现时间戳比较
		// if evt.Timestamp.After(sn.Timestamp()) {
		r.bus.Publish(evt)
		// }
	}
}
