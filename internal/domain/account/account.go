package account

import (
	"quant-trading/internal/domain/common"
	"quant-trading/internal/domain/user"

	"github.com/google/uuid"
)

/*
Account
=======

账户聚合根。

这里只保留基础身份信息。
*/

type AccountID string

func NewAccountID() AccountID {
	return AccountID(uuid.New().String())
}

func (a AccountID) String() string {
	return string(a)
}

type Status string

const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusSuspended Status = "suspended"
)

type Account struct {
	common.Model
	ID           AccountID      `gorm:"column:account_id;type:varchar(40);uniqueIndex;not null;primaryKey"`
	UserID       user.UserID    `gorm:"column:user_id;type:varchar(40);index;not null"`
	BrokerName   string         `gorm:"column:broker_name;not null"`
	AccountAlias string         `gorm:"column:account_alias;not null"`
	Status       Status         `gorm:"column:status;not null;default:active"`
	Config       map[string]any `gorm:"column:config;type:json"`
}

func (a *Account) TableName() string {
	return "t_account"
}

func NewAccount(userID user.UserID, brokerName, alias string) *Account {
	return &Account{
		ID:           NewAccountID(),
		UserID:       userID,
		BrokerName:   brokerName,
		AccountAlias: alias,
		Status:       StatusActive,
		Config:       make(map[string]any),
	}
}
