package instrument

type Type string

const (
	Stock  Type = "STOCK"  // 股票
	Future Type = "FUTURE" // 期货
	Option Type = "OPTION" // 期权
)

type OptionType string

const (
	Call OptionType = "CALL"
	Put  OptionType = "PUT"
)
