package backtest

import (
	"quant-trading/internal/domain/market"
	"time"
)

/*
DataSource
==========

数据源接口定义了回测引擎如何获取历史数据。

设计原则:
- 抽象数据来源(CSV/数据库/API)
- 支持多种数据格式
- 按时间顺序返回数据
*/
type DataSource interface {
	// GetEvents 获取指定时间的所有行情事件
	GetEvents(t time.Time) ([]market.Event, error)

	// GetRange 获取时间范围内的所有事件
	GetRange(start, end time.Time) ([]market.Event, error)
}

/*
MemoryDataSource
================

内存数据源实现,用于测试和简单回测。
*/
type MemoryDataSource struct {
	events []market.Event
	index  map[time.Time][]market.Event
}

// NewMemoryDataSource 创建内存数据源
func NewMemoryDataSource(events []market.Event) *MemoryDataSource {
	ds := &MemoryDataSource{
		events: events,
		index:  make(map[time.Time][]market.Event),
	}

	// 构建时间索引
	for _, evt := range events {
		// 按分钟对齐
		t := evt.Time.Truncate(time.Minute)
		ds.index[t] = append(ds.index[t], evt)
	}

	return ds
}

// GetEvents 获取指定时间的所有行情事件
func (ds *MemoryDataSource) GetEvents(t time.Time) ([]market.Event, error) {
	// 按分钟对齐
	t = t.Truncate(time.Minute)
	events, exists := ds.index[t]
	if !exists {
		return []market.Event{}, nil
	}
	return events, nil
}

// GetRange 获取时间范围内的所有事件
func (ds *MemoryDataSource) GetRange(start, end time.Time) ([]market.Event, error) {
	result := make([]market.Event, 0)
	for _, evt := range ds.events {
		if evt.Time.After(start) && evt.Time.Before(end) {
			result = append(result, evt)
		}
	}
	return result, nil
}
