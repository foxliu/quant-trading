package snapshot

type Store interface {
	Save(snapshot Snapshot)
	LoadLatest(name string) Snapshot
}
