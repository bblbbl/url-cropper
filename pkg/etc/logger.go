package etc

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger  *zap.Logger
	sLogger *zap.SugaredLogger
)

func InitLogger() {
	zapCnf := zap.Config{
		Level:            zap.NewAtomicLevel(),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
	}

	logger = zap.Must(zapCnf.Build())
	sLogger = logger.Sugar()
}

func FlushLogger() {
	err := logger.Sync()
	if err != nil {
		return
	}

	err = sLogger.Sync()
	if err != nil {
		return
	}
}

func GetLogger() *zap.SugaredLogger {
	return sLogger
}
