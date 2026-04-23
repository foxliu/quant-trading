package db

import (
	dAccount "quant-trading/internal/domain/account"
	"quant-trading/internal/domain/user"
	"quant-trading/internal/infrastructure/logger"
	"quant-trading/internal/infrastructure/persistence"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var log = logger.Logger.With(zap.String("module", "db"))

var DB *gorm.DB

func InitSQLite(dsn string) *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("打开数据库错误", zap.Error(err))
		panic(err)
	}
	// AutoMigrate(开发阶段)
	err = DB.AutoMigrate(
		&user.User{}, &user.Variety{}, &user.UserVariety{},
		&dAccount.Account{}, &dAccount.Snapshot{},
		&persistence.TraderRuntimeProfile{},
	)
	if err != nil {
		log.Error("自动迁移错误", zap.Error(err))
		panic(err)
	}
	return DB
}
