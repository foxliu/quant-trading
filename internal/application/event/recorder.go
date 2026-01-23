package event

type Recorder interface {
	Record(evt *Envelope)
	Events() []*Envelope
}
