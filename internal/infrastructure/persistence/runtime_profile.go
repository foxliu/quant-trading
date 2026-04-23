package persistence

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StrategySpec struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

type TraderRuntimeProfile struct {
	ID uint `gorm:"primaryKey"`

	Key    string `gorm:"column:profile_key;type:varchar(64);uniqueIndex;not null"`
	Active bool   `gorm:"column:active;not null;default:true"`

	BrokerName  string `gorm:"column:broker_name;type:varchar(32);not null;default:'CTP'"`
	AccountID   string `gorm:"column:account_id;type:varchar(64);not null"`
	InitialCash float64 `gorm:"column:initial_cash;not null;default:100000"`

	CTPFrontAddr string `gorm:"column:ctp_front_addr;type:varchar(255);not null"`
	CTPBrokerID  string `gorm:"column:ctp_broker_id;type:varchar(64);not null"`
	CTPUserID    string `gorm:"column:ctp_user_id;type:varchar(64);not null"`
	CTPInvestorID string `gorm:"column:ctp_investor_id;type:varchar(64);not null"`
	CTPPassword   string `gorm:"column:ctp_password;type:varchar(255);not null"`

	StrategiesJSON string `gorm:"column:strategies_json;type:text;not null"`
}

func (TraderRuntimeProfile) TableName() string {
	return "t_trader_runtime_profile"
}

func (p *TraderRuntimeProfile) Strategies() ([]StrategySpec, error) {
	if p.StrategiesJSON == "" {
		return []StrategySpec{}, nil
	}
	var strategies []StrategySpec
	if err := json.Unmarshal([]byte(p.StrategiesJSON), &strategies); err != nil {
		return nil, err
	}
	return strategies, nil
}

func GetActiveTraderRuntimeProfile(db *gorm.DB) (*TraderRuntimeProfile, error) {
	var profile TraderRuntimeProfile
	if err := db.Where("active = ?", true).Order("id asc").First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func UpsertTraderRuntimeProfile(db *gorm.DB, profile *TraderRuntimeProfile) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "profile_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"active", "broker_name", "account_id", "initial_cash", "ctp_front_addr", "ctp_broker_id", "ctp_user_id", "ctp_investor_id", "ctp_password", "strategies_json"}),
	}).Create(profile).Error
}
