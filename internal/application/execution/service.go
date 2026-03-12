package execution

import (
	"quant-trading/internal/application/account"
	"quant-trading/internal/domain/execution"
)

type Adapter interface {
	Execute(order execution.Order) ([]execution.Fill, error)
}

type Service struct {
	adapter Adapter
	account *account.Context
}

func NewService(adapter Adapter, account *account.Context) *Service {
	return &Service{
		adapter: adapter,
		account: account,
	}
}

func (s *Service) Submit(order execution.Order) error {
	fills, err := s.adapter.Execute(order)
	if err != nil {
		return err
	}

	for _, f := range fills {
		s.applyFill(f)
	}
	return nil
}

func (s *Service) applyFill(fill execution.Fill) {
	s.account.ApplyFill(
		fill.Symbol,
		fill.Side,
		fill.Price,
		fill.Qty,
	)
}
