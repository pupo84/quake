package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new instance of a SugaredLogger with a JSON encoder and an atomic level set to Info.
// The logger is configured to encode timestamps in ISO8601 format and the key for the timestamp is "timestamp".
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
