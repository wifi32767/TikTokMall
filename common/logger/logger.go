// 设置这个库是为了提供一个统一的日志接口，方便以后替换日志库
package logger

import "log/slog"

var LogLevel = LevelError

var (
	LevelDebug = -4
	LevelInfo  = 0
	LevelWarn  = 4
	LevelError = 8
)

func Init(level int) {
	LogLevel = level
	slog.SetLogLoggerLevel(slog.Level(LogLevel))
}

func Debug(msg string) {
	slog.Debug(msg)
}

func Info(msg string) {
	slog.Info(msg)
}

func Warn(msg string) {
	slog.Warn(msg)
}

func Error(msg string) {
	slog.Error(msg)
}
