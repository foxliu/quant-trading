package risk

type Rule interface {
	Name() string
	Evaluate(ctx *Context) *Result
}
