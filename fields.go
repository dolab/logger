package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		as.fields[key] = value
	}
}

// Bool is shortcut for bool field option.
func Bool(key string, value bool) Attr {
	return func(as *attrs) {
		as.fields[key] = value
	}
}

// Err is shortcut for error field option.
// NOTE: It uses error for the key forced!
func Err(err error) Attr {
	return func(as *attrs) {
		as.fields["error"] = err.Error()
	}
}

// Any for any type field option.
func Any(key string, value any) Attr {
	return func(as *attrs) {
		as.fields[key] = value
	}
}

type attrs struct {
	format Formatter
	fields map[string]any
	stacks []byte
}

func (as *attrs) IsValid() bool {
	if as == nil {
		return false
	}

	return len(as.fields) > 0
}

func (as *attrs) String() string {
	var buf bytes.Buffer
	n := len(as.fields)
	for k, v := range as.fields {
		buf.WriteString(k)
		buf.WriteString("=")

		switch t := v.(type) {
		case string:
			buf.WriteString(t)
		case []byte:
			buf.Write(t)
		case error:
			buf.WriteString(t.Error())
		default:
			if str, ok := v.(fmt.Stringer); ok {
				buf.WriteString(str.String())
			} else {
				buf.WriteString(fmt.Sprintf("%v", v))
			}
		}

		n--
		if n > 0 {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

func (as *attrs) JSONString() string {
	b, _ := json.Marshal(as.fields)
	return string(b)
}
