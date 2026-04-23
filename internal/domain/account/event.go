package account

type AccountCreateEvent struct {
	AccountID AccountID `json:"account_id"`
}

func (e AccountCreateEvent) String() string {
	return "AccountCreatedEvent{AccountID: " + string(e.AccountID) + "}"
}

type AccountBalanceChangedEvent struct {
	AccountID AccountID `json:"account_id"`
	Snapshot  Snapshot  `json:"snapshot"`
}

func (e AccountBalanceChangedEvent) String() string {
	return "AccountBalanceChangedEvent{AccountID: " + string(e.AccountID) + "}"
}
