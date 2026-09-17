package utils

import (
	"os"
	"time"

	"github.com/p1shiA/ezLink/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(cfg *config.Config) (*zap.Logger, error) {
	var log *zap.Logger

	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	// base encoder
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "ts"
	encCfg.EncodeTime = func (t time.Time, enc zapcore.PrimitiveArrayEncoder)  {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
	}
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	// terminal
	consoleCfg := encCfg
	consoleCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)

	// file
	fileEncoder := zapcore.NewJSONEncoder(encCfg)

	// rotating
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename: cfg.Filename,
		MaxSize: cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge: cfg.MaxAge,
		Compress: cfg.Compress,
	})

	stdout := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, fileWriter, zapcore.InfoLevel),
	)

	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return log, nil
}
