package logm

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogOption 配置Logger
type LogOption = zap.Option

// Logger 将消息记录到配置的输出中
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func newLogger(rootLogger *zap.Logger, selector string, options ...LogOption) *Logger {
	log := rootLogger.WithOptions(zap.AddCallerSkip(1)).
		WithOptions(options...).
		Named(selector)
	return &Logger{log, log.Sugar()}
}

// NewLogger 返回标有选择器的新Logger
func NewLogger(selector string, options ...LogOption) *Logger {
	return newLogger(loadLogger().rootLogger, selector, options...)
}

// IsDebug 检查是否启用了Debug模式
func (l *Logger) IsDebug() bool {
	return l.logger.Check(zapcore.DebugLevel, "") != nil
}

// Sprint

// Debug 使用fmt.Sprint构造和记录消息
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Info 使用fmt.Sprint构造和记录消息
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Warn 使用fmt.Sprint构造和记录消息
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Error 使用fmt.Sprint构造和记录消息
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Fatal 使用fmt.Sprint构造和记录消息, 会调用os.Exit(1)
func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// Panic 使用fmt.Sprint构造和记录消息, panic()
func (l *Logger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}

// DPanic 使用fmt.Sprint构造和记录消息, 在开发模式中，panic()
func (l *Logger) DPanic(args ...interface{}) {
	l.sugar.DPanic(args...)
}

// Sprintf

// Debugf 使用fmt.Sprintf构造和记录消息
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Infof 使用fmt.Sprintf构造和记录消息
func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

// Warnf 使用fmt.Sprintf构造和记录消息
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

// Errorf 使用fmt.Sprintf构造和记录消息
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

// Fatalf 使用fmt.Sprintf构造和记录消息, 会调用os.Exit(1)
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(format, args...)
}

// Panicf 使用fmt.Sprintf构造和记录消息, panic()
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sugar.Panicf(format, args...)
}

// DPanicf 使用fmt.Sprintf构造和记录消息, 在开发模式中，panic()
func (l *Logger) DPanicf(format string, args ...interface{}) {
	l.sugar.DPanicf(format, args...)
}

// 基于上下文

// Debugw 记录一条带有其他上下文的消息，额外的内容以键值对的方式
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

// Infow 记录一条带有其他上下文的消息，额外的内容以键值对的方式
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Errorw 记录一条带有其他上下文的消息，额外的内容以键值对的方式
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Fatalw 记录一条带有其他上下文的消息，额外的内容以键值对的方式
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

// Panicw 记录一条带有其他上下文的消息，额外的内容以键值对的方式
func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.sugar.Panicw(msg, keysAndValues...)
}

// DPanicw 记录一条带有其他上下文的消息，额外的内容以键值对的方式
func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.sugar.DPanicw(msg, keysAndValues...)
}

// Recover stops a panicking goroutine and logs an Error.
func (l *Logger) Recover(msg string) {
	if r := recover(); r != nil {
		msg := fmt.Sprintf("%s. Recovering, but please report this.", msg)
		l.Error(msg, zap.Any("panic", r), zap.Stack("stack"))
	}
}

// Sync syncs the logger.
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// L returns an unnamed global logger.
func L() *Logger {
	return loadLogger().logger
}