package event

// Recorder 永远不参与业务判断

type RecordingBus struct {
	bus      Bus
	recorder Recorder
}

func NewRecordingBus(bus Bus, recorder Recorder) *RecordingBus {
	return &RecordingBus{
		bus:      bus,
		recorder: recorder,
	}
}

func (b *RecordingBus) Publish(evt *Envelope) {
	b.recorder.Record(evt)
	b.bus.Publish(evt)
}

func (b *RecordingBus) Subscribe(t Type, h Handler) {
	b.bus.Subscribe(t, h)
}
