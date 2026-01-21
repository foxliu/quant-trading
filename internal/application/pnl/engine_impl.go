package pnl

import "context"

type engine struct {
	ctx *Context
}

func NewEngine(ctx *Context) Engine {
	return &engine{ctx: ctx}
}

func (e *engine) Start(ctx context.Context) error {
	return nil
}

func (e *engine) Stop() error {
	return nil
}
