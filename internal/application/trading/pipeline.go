package trading

import "context"

/*
Pipeline
========

交易管线占位：上游策略/风控串联尚未接入本包，保留包边界以便后续实现。
当前保证工程可编译；调用方请使用 runtime.Scheduler 等现有路径。
*/
type Pipeline struct{}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) Start(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}
