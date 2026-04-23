package logger

import (
	"os"
	"quant-trading/internal/infrastructure/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger(cfg *config.TraderConfig) error {
	level := zapcore.InfoLevel
	//fmt.Printf("config: %+v", cfg)
	//err := level.Set(cfg.AppLog.Level)
	//if err != nil {
	//	panic(fmt.Sprintf("log level error: %v", err))
	//}
	var core zapcore.Core
	if cfg.Mode == "dev" {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeDuration = zapcore.MillisDurationEncoder
		encoderCfg.EncodeCaller = zapcore.FullCallerEncoder
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			level,
		)
	} else {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeDuration = zapcore.MillisDurationEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   cfg.Logging.Filename,
				MaxSize:    cfg.Logging.MaxSize,
				MaxAge:     cfg.Logging.MaxAge,
				MaxBackups: cfg.Logging.MaxBackups,
			}),
			level,
		)
	}
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	zap.ReplaceGlobals(logger)
	Logger = zap.L()
	return nil
}

func Sync() {
	_ = Logger.Sync()
}
