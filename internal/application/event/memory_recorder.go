package event

import "sync"

type MemoryRecorder struct {
	mu     sync.Mutex
	events []*Envelope
}

func NewMemoryRecorder() *MemoryRecorder {
	return &MemoryRecorder{
		events: make([]*Envelope, 0),
	}
}

func (r *MemoryRecorder) Record(evt *Envelope) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, evt)
}

func (r *MemoryRecorder) Events() []*Envelope {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]*Envelope(nil), r.events...)
}

func (r *MemoryRecorder) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = r.events[:0]
}
