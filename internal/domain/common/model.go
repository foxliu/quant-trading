package common

import (
	"database/sql"
	"time"
)

type Model struct {
	ID        uint         `gorm:"primary_key" json:"id"`
	CreatedAt time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time    `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"column:deleted_at;index" json:"deleted_at"`
}
