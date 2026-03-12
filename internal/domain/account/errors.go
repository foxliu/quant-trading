package account

import "errors"

var (
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrInsufficientPosition = errors.New("insufficient position")
)
