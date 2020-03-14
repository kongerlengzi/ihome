package zaplog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger


func Init()  {
	writerSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder,writerSyncer,zapcore.DebugLevel)
	sLogger := zap.New(core,zap.AddCaller())
	Logger = sLogger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:       "./test.log",
		MaxSize:        1,
		MaxBackups:     2,
		MaxAge:         30,
		Compress:       false,
	}
	return zapcore.AddSync(lumberjackLogger)
}

