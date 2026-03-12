package account

/*
Balance
=======

Balance 描述账户的资金状态
*/
type Balance struct {
	available float64
	frozen    float64
}

func NewBalance(initial float64) *Balance {
	return &Balance{available: initial}
}

func (b *Balance) Available() float64 {
	return b.available
}

func (b *Balance) Frozen() float64 {
	return b.frozen
}

func (b *Balance) Freeze(amount float64) {
	b.available -= amount
	b.frozen += amount
}

func (b *Balance) Unfreeze(amount float64) {
	b.frozen -= amount
	b.available += amount
}

func (b *Balance) Deduct(amount float64) {
	b.available -= amount
}

func (b *Balance) Add(amount float64) {
	b.available += amount
}

func (b *Balance) Snapshot() BalanceSnapshot {
	return BalanceSnapshot{
		Available: b.available,
		Frozen:    b.frozen,
	}
}

func (b *Balance) Restore(s BalanceSnapshot) {
	b.available = s.Available
	b.frozen = s.Frozen
}
