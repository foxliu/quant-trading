package account

import (
	"fmt"
	aCapital "quant-trading/internal/application/capital"
	"quant-trading/internal/application/event"
	"quant-trading/internal/domain/account"
	dCapital "quant-trading/internal/domain/capital"
	"quant-trading/internal/domain/portfolio"
	"quant-trading/internal/domain/user"
	"quant-trading/internal/infrastructure/logger"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	repo account.Repository
	bus  event.Bus
	log  *zap.Logger
}

func NewService(repo account.Repository, bus event.Bus) *Service {
	return &Service{
		repo: repo,
		bus:  bus,
		log:  logger.Logger.With(zap.String("module", "account.service"))}
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
		Type:      event.EventAccountCreated,
		Source:    "account.service",
		Timestamp: time.Now(),
		Payload:   account.AccountCreateEvent{AccountID: acc.ID},
	})
	return acc, nil
}

func (s *Service) NewContext(id account.AccountID, cap dCapital.Engine, port portfolio.Engine) (*Context, error) {
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

func (s *Service) NewCapitalEngineForAccount(acc *account.Account) dCapital.Engine {
	initialCapital := 10000.0 // 默认值
	if val, ok := acc.Config["initial_capital"]; ok {
		initialCapital = val.(float64)
	}

	// 从最新的 Snapshot 读取
	snap, _ := s.repo.GetLatestSnapshot(acc.ID)
	if snap != nil {
		if snap.Balance.Available > 0 {
			initialCapital = snap.Balance.Available
		}
	}
	return aCapital.NewMemoryEngine(initialCapital)
}

type LoadAccountCommand struct {
	UserID     user.UserID
	BrokerName string
	Alias      string
	Config     map[string]any
}

func (s *Service) LoadOrCreateAccount(cmd LoadAccountCommand) (*account.Account, error) {
	stableID := account.AccountID(cmd.Alias)
	if stableID == "" {
		return nil, fmt.Errorf("account: alias is required and used as stable account_id")
	}

	existing, err := s.repo.FindByID(stableID)
	if err == nil {
		return existing, nil
	}
	s.log.Info("按稳定账户ID加载失败，尝试创建账户", zap.String("account_id", stableID.String()), zap.Error(err))

	s.log.Warn("账户不存在，创建新账户", zap.Any("account", cmd))
	acc := account.NewAccount(cmd.UserID, cmd.BrokerName, cmd.Alias)
	acc.ID = stableID
	if cmd.Config != nil {
		acc.Config = cmd.Config
	} else {
		acc.Config = make(map[string]any)
	}
	if err := s.repo.SaveAccount(acc); err != nil {
		s.log.Error("保存账户失败", zap.Error(err), zap.Any("account", acc))
		return nil, err
	}

	s.bus.Publish(&event.Envelope{
		Type:      event.EventAccountCreated,
		Source:    "account.service",
		Timestamp: time.Now(),
		Payload:   account.AccountCreateEvent{AccountID: acc.ID},
	})

	return acc, nil
}
