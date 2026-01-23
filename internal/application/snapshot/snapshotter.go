package snapshot

/*
Snapshotter 是 Context 的能力

Event 系统不关心 Snapshot 内容
*/

type Snapshotter interface {
	Take() Snapshot
	Restore(snapshot Snapshot) error
}
