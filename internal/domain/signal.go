package domain

import "time"

/*
最小定义
*/

type Side string

const (
	Buy  Side = "BUY"
	Sell Side = "SELL"
)

type Signal struct {
	Instrument string
	Time       time.Time
	Side       Side
	Strength   float64
	Reason     string
}
