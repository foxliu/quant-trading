package account

import "quant-trading/internal/domain/user"

type Repository interface {
	SaveAccount(a *Account) error
	FindByID(id AccountID) (*Account, error)
	ListByUser(userID user.UserID) ([]Account, error)
	SaveSnapshot(s *Snapshot) error
	GetLatestSnapshot(accountID AccountID) (*Snapshot, error)
}
