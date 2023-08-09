package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.SugaredLogger {
	atomic := zap.NewAtomicLevel()

	atomic.SetLevel(zapcore.InfoLevel)

	enconfig := zap.NewProductionEncoderConfig()
	enconfig.EncodeTime = zapcore.ISO8601TimeEncoder
	enconfig.TimeKey = "timestamp"

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(enconfig),
		zapcore.Lock(os.Stdout),
		atomic,
	))

	defer logger.Sync()

	return logger.Sugar()
}
