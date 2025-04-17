package logger

import (
	"context"
	"log/slog"
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// Log drops in of slog.Log
func (l *Logger) Log(ctx context.Context, slogLevel slog.Level, msg string, args ...any) {
	log := l.NewTextLogger()
	switch slogLevel {
	case slog.LevelDebug:
		log.Debugf(msg, args...)
	case slog.LevelInfo:
		log.Infof(msg, args...)
	case slog.LevelWarn:
		log.Warnf(msg, args...)
	case slog.LevelError:
		log.Errorf(msg, args...)
	default:
		log.Infof(msg, args...)
	}
}

// LogAttrs drops in of slog.LogAttrs
func (l *Logger) LogAttrs(ctx context.Context, slogLevel slog.Level, msg string, slogAttrs ...slog.Attr) {
	var level = Ltrace
	switch slogLevel {
	case slog.LevelDebug:
		level = Ldebug
	case slog.LevelInfo:
		level = Linfo
	case slog.LevelWarn:
		level = Lwarn
	case slog.LevelError:
		level = Lerror
	}

	log := l.NewTextLogger().(*structLog)
	_ = log.writer(level, &attrs{
		format: log.format,
		stacks: log.stacks,
		fields: slogAttrs,
	}, msg)
}
