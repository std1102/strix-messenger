package system

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

var Logger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	return Logger
}

func InitLog() {
	Logger = initZapLog()
}

func initZapLog() *zap.SugaredLogger {
	cfg := zap.Config{
		Encoding:    "console",
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			TimeKey:      "time",
			LevelKey:     "level",
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeLevel:  customLevelEncoder,
			EncodeTime:   syslogTimeEncoder,
		},
	}

	logger, _ := cfg.Build()
	return logger.Sugar()
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	level, exist := loggerLevelMap[SystemConfig.Log.Level]
	if exist {
		enc.AppendString("[" + level.CapitalString() + "]")
	} else {
		enc.AppendString("[" + loggerLevelMap["debug"].CapitalString() + "]")
	}
}
