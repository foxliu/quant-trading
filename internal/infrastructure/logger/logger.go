package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

var cfg = struct {
	AppLog struct {
		Filename   string `yaml:"filename"`
		MaxSize    int    `yaml:"max_size"`
		MaxAge     int    `yaml:"max_age"`
		MaxBackups int    `yaml:"max_backups"`
	} `yaml:"app_log"`
	Mode string `yaml:"mode"`
}{
	AppLog: struct {
		Filename   string `yaml:"filename"`
		MaxSize    int    `yaml:"max_size"`
		MaxAge     int    `yaml:"max_age"`
		MaxBackups int    `yaml:"max_backups"`
	}{
		Filename:   "app.log",
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 3,
	},
	Mode: "dev",
}

func InitLogger() error {
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
				Filename:   cfg.AppLog.Filename,
				MaxSize:    int(cfg.AppLog.MaxSize),
				MaxAge:     int(cfg.AppLog.MaxAge),
				MaxBackups: int(cfg.AppLog.MaxBackups),
			}),
			level,
		)
	}
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	zap.ReplaceGlobals(logger)
	return nil
}

func Sync() {
	_ = Logger.Sync()
}
