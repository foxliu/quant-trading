package account

import (
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/account"
	"quant-trading/internal/domain/capital"
	"quant-trading/internal/domain/portfolio"
	"quant-trading/internal/domain/user"
	"quant-trading/internal/infrastructure/logger"

	"go.uber.org/zap"
)

type Service struct {
	repo account.Repository
	bus  event.Bus
	log  *zap.Logger
}

func NewService(repo account.Repository, bus event.Bus) *Service {
	return &Service{repo: repo, bus: bus, log: logger.Logger.With(zap.String("module", "account.service"))}
}

type CreateAccountCommand struct {
	UserID     user.UserID
	BrokerName string
	Alias      string
}

func (s *Service) CreateAccount(cmd CreateAccountCommand) (*account.Account, error) {
	acc := account.NewAccount(cmd.UserID, cmd.BrokerName, cmd.Alias)
	if err := s.repo.SaveAccount(acc); err != nil {
		s.log.Error("保存账户失败", zap.Error(err), zap.Any("account", acc))
		return nil, err
	}
	s.bus.Publish(&event.Envelope{
		Type: event.EventAccountCreated,
	})
	return acc, nil
}

func (s *Service) NewContext(id account.AccountID, cap capital.Engine, port portfolio.Engine) (*Context, error) {
	acc, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	snap, _ := s.repo.GetLatestSnapshot(id)

	ctx := NewContext(acc, s.bus, cap, port)

	if snap != nil {
		ctx.Restore(*snap)
	}
	return ctx, nil
}
