package logger

import (
	"fmt"
	"runtime/debug"
	"time"
)

// StructLogger for well formatted
type (
	StructLogger interface {
		Str(key, value string) StructLogger
		Bool(key string, value bool) StructLogger
		Duration(key string, value time.Duration) StructLogger
		Time(key string, value time.Time) StructLogger
		Err(err error, stack bool) StructLogger
		Any(key string, value any) StructLogger
		Fields(fields map[string]any) StructLogger

		Debug(msg string)
		Debugf(format string, args ...any)
		Info(msg string)
		Infof(format string, args ...any)
		Warn(msg string)
		Warnf(format string, args ...any)
		Error(msg string)
		Errorf(format string, args ...any)
		Fatal(msg string)
		Fatalf(format string, args ...any)
		Panic(msg string)
		Panicf(format string, args ...any)
	}
)

var (
	_ StructLogger = (*structLog)(nil)
)

type structLog struct {
	writer func(level Level, as *attrs, msg string) error
	format Formatter
	attrs  []Attr
	stacks []byte
}

func (log *structLog) Str(key, value string) StructLogger {
	log.attrs = append(log.attrs, String(key, value))
	return log
}

func (log *structLog) Bool(key string, value bool) StructLogger {
	log.attrs = append(log.attrs, Bool(key, value))
	return log
}

func (log *structLog) Duration(key string, value time.Duration) StructLogger {
	log.attrs = append(log.attrs, Duration(key, value))
	return log
}

func (log *structLog) Time(key string, value time.Time) StructLogger {
	log.attrs = append(log.attrs, Time(key, value))
	return log
}

func (log *structLog) Any(key string, value any) StructLogger {
	log.attrs = append(log.attrs, Any(key, value))
	return log
}

func (log *structLog) Err(err error, stack bool) StructLogger {
	if err == nil {
		return log
	}

	log.attrs = append(log.attrs, Err(err))
	if stack {
		log.stacks = debug.Stack()
	}
	return log
}

func (log *structLog) Fields(fields map[string]any) StructLogger {
	for k, v := range fields {
		switch t := v.(type) {
		case string:
			log.Str(k, t)
		case []byte:
			log.Str(k, string(t))
		case bool:
			log.Bool(k, t)
		case time.Duration:
			log.Duration(k, t)
		case time.Time:
			log.Time(k, t)
		case error:
			log.Err(t, k == "true")
		default:
			log.Any(k, t)
		}
	}

	return log
}

func (log *structLog) fields() *attrs {
	as := &attrs{
		format: log.format,
		stacks: log.stacks,
	}
	for _, attr := range log.attrs {
		attr(as)
	}

	return as
}

func (log *structLog) Debug(msg string) {
	_ = log.writer(Ldebug, log.fields(), msg)
}

func (log *structLog) Debugf(format string, args ...any) {
	_ = log.writer(Ldebug, log.fields(), fmt.Sprintf(format, args...))
}

func (log *structLog) Info(msg string) {
	_ = log.writer(Linfo, log.fields(), msg)
}

func (log *structLog) Infof(format string, args ...any) {
	_ = log.writer(Linfo, log.fields(), fmt.Sprintf(format, args...))
}

func (log *structLog) Warn(msg string) {
	_ = log.writer(Lwarn, log.fields(), msg)
}

func (log *structLog) Warnf(format string, args ...any) {
	_ = log.writer(Lwarn, log.fields(), fmt.Sprintf(format, args...))
}

func (log *structLog) Error(msg string) {
	_ = log.writer(Lerror, log.fields(), msg)
}

func (log *structLog) Errorf(format string, args ...any) {
	_ = log.writer(Lerror, log.fields(), fmt.Sprintf(format, args...))
}

func (log *structLog) Fatal(msg string) {
	_ = log.writer(Lfatal, log.fields(), msg)
}

func (log *structLog) Fatalf(format string, args ...any) {
	_ = log.writer(Lfatal, log.fields(), fmt.Sprintf(format, args...))
}

func (log *structLog) Panic(msg string) {
	_ = log.writer(Lpanic, log.fields(), msg)
}

func (log *structLog) Panicf(format string, args ...any) {
	_ = log.writer(Lpanic, log.fields(), fmt.Sprintf(format, args...))
}
