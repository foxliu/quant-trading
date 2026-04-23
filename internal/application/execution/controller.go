package execution

type Controller interface {
	Execute(cmd Command) error
}
