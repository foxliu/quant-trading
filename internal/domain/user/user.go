package user

import (
	"quant-trading/internal/domain/common"
)

type UserID string

func (u UserID) String() string {
	return string(u)
}

type User struct {
	common.Model
	ID           UserID         `gorm:"column:user_id;type:varchar(40);uniqueIndex;not null;primaryKey"`
	BrokerID     string         `gorm:"column:broker_id;not null;comment:证卷公司ID"`
	InvestorID   string         `gorm:"column:investor_id;not null;comment:投资者代码，交易账号"`
	AppID        string         `gorm:"column:app_id;not null;comment:穿透式监管客户端应用标识"`
	AuthCode     string         `gorm:"column:auth_code;not null;comment:穿透式监管授权码/认证码"`
	ServerJson   []byte         `gorm:"column:server_json;type:JSONB;not null"`
	PasswordHash string         `gorm:"column:password_hash;not null"`
	IsActive     bool           `gorm:"column:is_active;not null"`
}

func (u *User) TableName() string {
	return "t_user"
}

/*
全局品种库，用于记录交易的品种信息，如：A股、B股、期货、期权
*/

type Variety struct {
	common.Model
	Code        string `gorm:"column:code" json:"code"`
	Name        string `gorm:"column:name" json:"name"`
	Market      string `gorm:"column:market" json:"market"`
	Exchange    string `gorm:"column:exchange" json:"exchange"`
	Description string `gorm:"column:description" json:"description"`
}

func (v *Variety) TableName() string {
	return "t_variety"
}

type UserVariety struct {
	common.Model
	UserID    uint    `gorm:"column:user_id" json:"user_id"`
	User      User    `gorm:"-" json:"user"`
	VarietyID uint    `gorm:"column:variety_id" json:"variety_id"`
	Variety   Variety `gorm:"-" json:"variety"`
}

func (v *UserVariety) TableName() string {
	return "t_user_variety"
}
