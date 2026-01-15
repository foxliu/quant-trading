package order

/*
Type
====

Order Type 描述订单的执行方式。
*/
type Type string

const (
	Market Type = "MARKET"
	Limit  Type = "LIMIT"
)
