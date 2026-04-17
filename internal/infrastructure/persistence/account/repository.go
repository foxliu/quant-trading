package account

import (
	"quant-trading/internal/domain/account"
	"quant-trading/internal/domain/user"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) SaveAccount(a *account.Account) error {
	return r.db.Save(a).Error
}

func (r *repository) FindByID(id account.AccountID) (*account.Account, error) {
	var a account.Account
	err := r.db.First(&a, "account_id = ?", id).Error
	return &a, err
}

func (r *repository) ListByUser(userID user.UserID) ([]account.Account, error) {
	var accounts []account.Account
	err := r.db.Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

func (r *repository) SaveSnapshot(s *account.Snapshot) error {
	return r.db.Save(s).Error
}

func (r *repository) GetLatestSnapshot(accountID account.AccountID) (*account.Snapshot, error) {
	var s account.Snapshot
	err := r.db.Where("account_id = ?", accountID).Order("timestamp desc").First(&s).Error
	return &s, err
}

func NewRepository(db *gorm.DB) account.Repository {
	return &repository{db: db}
}
