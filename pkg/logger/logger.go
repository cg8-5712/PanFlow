package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(level string) error {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	log, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}

func Info(msg string, fields ...zap.Field)   { log.Info(msg, fields...) }
func Infof(format string, args ...interface{}) { log.Sugar().Infof(format, args...) }
func Warn(msg string, fields ...zap.Field)   { log.Warn(msg, fields...) }
func Warnf(format string, args ...interface{}) { log.Sugar().Warnf(format, args...) }
func Error(msg string, fields ...zap.Field)  { log.Error(msg, fields...) }
func Errorf(format string, args ...interface{}) { log.Sugar().Errorf(format, args...) }
func Debug(msg string, fields ...zap.Field)  { log.Debug(msg, fields...) }
func Debugf(format string, args ...interface{}) { log.Sugar().Debugf(format, args...) }
func Fatal(msg string, fields ...zap.Field)  { log.Fatal(msg, fields...) }
func Fatalf(format string, args ...interface{}) { log.Sugar().Fatalf(format, args...) }
