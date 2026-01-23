package snapshot

import "sync"

type MemoryStore struct {
	mu        sync.Mutex
	snapshots map[string]Snapshot
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		snapshots: make(map[string]Snapshot),
	}
}

func (s *MemoryStore) Save(sn Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.snapshots[sn.Name()] = sn
}

func (s *MemoryStore) LoadLatest(name string) Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshots[name]
}
