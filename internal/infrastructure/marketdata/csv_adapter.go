package marketdata

import (
	"encoding/csv"
	"fmt"
	"os"
	"quant-trading/internal/application/backtest"
	"quant-trading/internal/domain/instrument"
	"quant-trading/internal/domain/market"
	"strconv"
	"time"
)

/*
CSVAdapter
==========

CSV 文件行情数据适配器（基础设施层）。

支持标准 OHLCV CSV 格式：
timestamp,open,high,low,close,volume
2025-03-01 09:30:00,100.50,101.20,100.10,100.80,12500

设计原则:
- 预加载全部数据到内存（复用 MemoryDataSource）
- 自动构造 market.Bar + market.Event
- 支持任意 Instrument（股票/期货/期权）
- 错误鲁棒：跳过空行、解析失败行仅警告
*/
type CSVAdapter struct {
	ds *backtest.MemoryDataSource
}

// NewCSVAdapter 创建 CSV 数据适配器
// filePath: CSV 文件绝对/相对路径
// instr:    完整合约模型（Symbol、Exchange、Type 等）
func NewCSVAdapter(filePath string, instr instrument.Instrument) (*CSVAdapter, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开 CSV 文件失败: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // 允许变长列

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取 CSV 文件失败: %w", err)
	}

	events := make([]market.Event, 0, len(records))

	for i, row := range records {
		// 跳过空行或注释行
		if len(row) == 0 || row[0] == "" || row[0][0] == '#' {
			continue
		}
		// 第一行可能是header, 自动跳过（timestamp 开头）
		if i == 0 && (row[0] == "timestamp" || row[0] == "time" || row[0] == "date") {
			continue
		}

		if len(row) < 6 {
			fmt.Printf("[WARN] CSV第 %d 行列数不足，跳过\n", i+1)
			continue
		}

		// 解析时间（支持两种常见格式）
		t, err := parseTime(row[0])
		if err != nil {
			fmt.Printf("[WARN] CSV第 %d 行时间解析失败: %v，跳过\n", i+1, err)
			continue
		}

		open, _ := strconv.ParseFloat(row[1], 64)
		high, _ := strconv.ParseFloat(row[2], 64)
		low, _ := strconv.ParseFloat(row[3], 64)
		closePrice, _ := strconv.ParseFloat(row[4], 64)
		volume, _ := strconv.ParseFloat(row[5], 64)

		bar := market.Bar{
			Instrument: instr,
			Time:       t,
			Open:       open,
			High:       high,
			Low:        low,
			Close:      closePrice,
			Volume:     volume,
		}

		evt := market.Event{
			Type: market.EventMarket,
			Time: t,
			Data: bar,
		}

		events = append(events, evt)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("CSV 中未解析到任何有效行情数据")
	}

	// 委托给已有 MemoryDataSource（复用 GetEvents/GetRange）
	ds := backtest.NewMemoryDataSource(events)
	return &CSVAdapter{ds: ds}, nil
}

// GetEvents 实现backtest.DataSource接口
func (a *CSVAdapter) GetEvents(t time.Time) ([]market.Event, error) {
	return a.ds.GetEvents(t)
}

// GetRange 实现backtest.DataSource接口
func (a *CSVAdapter) GetRange(start, end time.Time) ([]market.Event, error) {
	return a.ds.GetRange(start, end)
}

// parseTime 解析时间
func parseTime(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
		"2006.01.02 15:04:05",
		"2006.01.02 15:04",
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析时间: %s", s)
}
