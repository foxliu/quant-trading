package db

import (
	"quant-trading/internal/domain/user"
	"quant-trading/internal/infrastructure/logger"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var log = logger.Logger.With(zap.String("module", "db"))

var DB *gorm.DB

func InitSQLite(dsn string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("打开数据库错误", zap.Error(err))
		panic(err)
	}
	// AutoMigrate(开发阶段)
	err = DB.AutoMigrate(&user.User{}, &user.Variety{}, &user.UserVariety{})
	if err != nil {
		log.Error("自动迁移错误", zap.Error(err))
		panic(err)
	}
}
