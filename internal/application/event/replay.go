package event

/*
Replay 的核心思想

Replay 不是“重放函数”，而是“重新驱动整个系统”
*/

type Replayer struct {
	bus Bus
}

func NewReplayer(bus Bus) *Replayer {
	return &Replayer{bus: bus}
}

func (r *Replayer) Replay(events []*Envelope) {
	for _, evt := range events {
		r.bus.Publish(evt)
	}
}
