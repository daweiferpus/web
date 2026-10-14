package public

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() *zap.SugaredLogger {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewTee(zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel), // 写入文件
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)) // 写入控制台

	logger := zap.New(core, zap.AddCaller())
	LOG = logger.Sugar()
	return LOG
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log/rkstk.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
