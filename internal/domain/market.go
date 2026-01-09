package domain

import "time"

/*
最小定义
*/

type MarketBar struct {
	Instrument string
	Time       time.Time
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Volume     float64
}
