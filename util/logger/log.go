package logger

import (
	"log"

	"go.uber.org/zap"
)

var zlog *zap.SugaredLogger

//todo , if no path ,change to init func
func InitLog() {
	logger, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal("init zap logger failed: ", err.Error())
	}
	zlog = logger.Sugar()
}

func Info(args ...interface{}) {
	zlog.Info(args...)
}

func Infof(template string, args ...interface{}) {
	zlog.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	zlog.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	zlog.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	zlog.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	zlog.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	zlog.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	zlog.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	zlog.Errorw(msg, keysAndValues...)
}
