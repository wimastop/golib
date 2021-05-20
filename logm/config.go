package logm

import "time"

// Config 日志的配置项
type Config struct {
	JSON      bool     // 日志输出格式为JSON
	Level     Level    // 日志等级 (error, warning, info, debug)
	Selectors []string // 用于调试级别日志记录的选择器

	Files FileConfig

	environment Environment

	addCaller   bool
	development bool
}

type FileConfig struct {
	Path            string
	Name            string
	MaxSize         uint
	MaxBackups      uint
	Permissions     uint32
	Interval        time.Duration
	RotateOnStartup bool
	RedirectStderr  bool
}

const (
	defaultLevel = InfoLevel
)

func DefaultConfig(environment Environment) Config {
	return Config{
		Level: defaultLevel,
		Files: FileConfig{
			MaxSize:         10 * 1024 * 1024,
			MaxBackups:      7,
			Permissions:     0600,
			Interval:        0,
			RotateOnStartup: true,
		},
		environment: environment,
		addCaller:   true,
	}
}
