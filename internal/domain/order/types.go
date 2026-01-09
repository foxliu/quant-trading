package order

type Type string

const (
	Market Type = "MARKET"
	Limit  Type = "LIMIT"
)

type Status string

const (
	New       Status = "NEW"
	Submitted Status = "SUBMITTED"
	Filled    Status = "FILLED"
	Canceled  Status = "CANCELED"
	Rejected  Status = "REJECTED"
)
