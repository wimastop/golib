package logm

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"strings"
)

// Level 日志记录优先级，级别越高越重要
type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

var levelStrings = map[Level]string{
	DebugLevel: "debug",
	InfoLevel:  "info",
	WarnLevel:  "warning",
	ErrorLevel: "error",
}

var zapLevels = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
}

// String 返回日志级别的名称
func (l Level) String() string {
	s, found := levelStrings[l]
	if found {
		return s
	}
	return fmt.Sprintf("Level(%d)", l)
}

// Enabled 返回是否启用给定的级别
func (l Level) Enabled(level Level) bool {
	return level >= l
}

// ZapLevel 返回Zap的日志等级
func (l Level) ZapLevel() zapcore.Level {
	z, found := zapLevels[l]
	if found {
		return z
	}
	return zapcore.InfoLevel
}

// Unpack 解析string类型为日志等级
func (l *Level) Unpack(str string) error {
	str = strings.ToLower(str)

	for level, name := range levelStrings {
		if name == str {
			*l = level
			return nil
		}
	}

	return errors.Errorf("invalid level '%v'", str)
}

// MarshalYAML marshals 日志等级
func (l Level) MarshalYAML() (interface{}, error) {
	s, found := levelStrings[l]
	if found {
		return s, nil
	}
	return nil, errors.Errorf("invalid level '%d'", l)
}

// MarshalJSON marshals 日志等级
func (l Level) MarshalJSON() ([]byte, error) {
	s, found := levelStrings[l]
	if found {
		return []byte(s), nil
	}
	return nil, errors.Errorf("invalid level '%d'", l)
}
