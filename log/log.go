package log

import (
	"github.com/zsmartex/pkg/services"
)

var Logger = services.NewLoggerService("Kouda")

func Info(args ...interface{}) {
	Logger.Info(args)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args)
}

func Debug(args ...interface{}) {
	Logger.Debug(args)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args)
}

func Error(args ...interface{}) {
	Logger.Error(args)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args)
}

func Panic(args ...interface{}) {
	Logger.Panic(args)
}

func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args)
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args)
}

func Warn(args ...interface{}) {
	Logger.Warn(args)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args)
}
