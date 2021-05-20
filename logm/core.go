package logm

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	golog "log"
	"sync/atomic"
	"unsafe"
)

var (
	_log          unsafe.Pointer // 指向coreLogger, 通过atomic.LoadPointer访问
	_defaultGoLog = golog.Writer()
)

func init() {
	storeLogger(&coreLogger{
		selectors:    map[string]struct{}{},
		rootLogger:   zap.NewNop(),
		globalLogger: zap.NewNop(),
		logger:       newLogger(zap.NewNop(), ""),
	})
}

type coreLogger struct {
	selectors    map[string]struct{}    // 启用的选择器
	rootLogger   *zap.Logger            // 根记录器，未配置任何选项
	globalLogger *zap.Logger            // 全局功能使用的记录器
	logger       *Logger                // 所有logm.Logger的基础记录器
	observedLogs *observer.ObservedLogs // Contains events generated while in observation mode (a testing mode).
}

func globalLogger() *zap.Logger {
	return loadLogger().globalLogger
}

func loadLogger() *coreLogger {
	p := atomic.LoadPointer(&_log)
	return (*coreLogger)(p)
}

func storeLogger(l *coreLogger) {
	if old := loadLogger(); old != nil {
		_ = old.rootLogger.Sync()
	}
	atomic.StorePointer(&_log, unsafe.Pointer(l))
}
