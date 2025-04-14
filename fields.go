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
// NOTE: It translate bool to string, such as true="true", false="false".
func Bool(key string, value bool) Attr {
	return func(as *attrs) {
		if value {
			as.fields[key] = "true"
		} else {
			as.fields[key] = "false"
		}
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
}

func (as attrs) String() string {
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

func (as attrs) JSONString() string {
	b, _ := json.Marshal(as.fields)
	return string(b)
}
