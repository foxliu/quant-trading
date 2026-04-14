package user_model

import (
	"database/sql"
	"quant-trading/internal/infrastructure/logger"
	"time"

	"go.uber.org/zap"
)

var log = logger.Logger.With(zap.String("module", "models"))

type Model struct {
	ID        uint         `gorm:"primary_key" json:"id"`
	CreatedAt time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time    `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"column:deleted_at;index" json:"deleted_at"`
}
