package account

import (
	"context"
	"time"
)

type RefreshScheduler interface {
	Start(ctx context.Context) error // 启动定时刷新
	Stop() error                     // 优雅停止
	SetInterval(d time.Duration)     // 动态调整间隔
}
