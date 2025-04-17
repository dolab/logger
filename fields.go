package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

const (
	TextFormat Formatter = iota
	JSONFormat
)

type (
	// Attr for fields option
	Attr func(as *attrs)

	Formatter int
)

// String is shortcut for string field option.
func String(key, value string) Attr {
	return func(as *attrs) {
		as.fields = append(as.fields, slog.String(key, value))
	}
}

// Bool is shortcut for bool field option.
func Bool(key string, value bool) Attr {
	return func(as *attrs) {
		as.fields = append(as.fields, slog.Bool(key, value))
	}
}

// Duration is shortcut for time.Duration field option.
func Duration(key string, value time.Duration) Attr {
	return func(as *attrs) {
		as.fields = append(as.fields, slog.String(key, fmt.Sprintf("%v", value)))
	}
}

// Time is shortcut for time.Time field option.
func Time(key string, value time.Time) Attr {
	return func(as *attrs) {
		as.fields = append(as.fields, slog.Time(key, value))
	}
}

// Err is shortcut for error field option.
// NOTE: It uses error for the key forced!
func Err(err error) Attr {
	return func(as *attrs) {
		as.fields = append(as.fields, slog.Any("error", err.Error()))
	}
}

// Any for any type field option.
func Any(key string, value any) Attr {
	return func(as *attrs) {
		as.fields = append(as.fields, slog.Any(key, value))
	}
}

type attrs struct {
	format Formatter
	fields []slog.Attr
	stacks []byte
}

func (as *attrs) IsValid() bool {
	if as == nil {
		return false
	}

	return len(as.fields) > 0
}

func (as *attrs) String() string {
	n := len(as.fields) - 1

	var buf bytes.Buffer
	for i, attr := range as.fields {
		buf.WriteString(attr.Key)
		buf.WriteString("=")
		buf.WriteString(attr.Value.String())

		if n > i {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

func (as *attrs) JSONString() string {
	fields := make(map[string]any)
	for _, attr := range as.fields {
		switch attr.Value.Kind() {
		case slog.KindString:
			fields[attr.Key] = attr.Value.String()
		case slog.KindBool:
			fields[attr.Key] = attr.Value.Bool()
		case slog.KindDuration:
			fields[attr.Key] = fmt.Sprintf("%v", attr.Value.Duration())
		case slog.KindTime:
			fields[attr.Key] = attr.Value.Time()
		case slog.KindAny:
			fields[attr.Key] = attr.Value
		default:
			fields[attr.Key] = attr.Value.String()
		}
	}

	b, _ := json.Marshal(fields)
	return string(b)
}
