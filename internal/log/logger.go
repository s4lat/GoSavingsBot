package log

import (
	"go.uber.org/zap"
)

var sl *zap.SugaredLogger
var l *zap.Logger

func init() {
	l, _ = zap.NewProduction()

	sl = l.Sugar()
	sl.Info("logger initialized")
}

//func Logger() *zap.Logger {
//	return l
//}

func Sugar() *zap.SugaredLogger {
	return sl
}

func Close() error {
	return l.Sync() // flushes buffer
}
