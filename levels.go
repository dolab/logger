package logger

import (
	"strings"
)

const (
	lmin Level = iota
	Ldebug
	Linfo
	Lwarn
	Lerror
	Lfatal
	Lpanic
	Ltrace
	lmax
)

var (
	// Logger levels
	levels = map[Level]string{
		Ldebug: "DEBUG",
		Linfo:  "INFO",
		Lwarn:  "WARN",
		Lerror: "ERROR",
		Lfatal: "FATAL",
		Lpanic: "PANIC",
		Ltrace: "Stack",
	}
)

type Level int

func (l Level) IsValid() bool {
	return lmin < l && l < lmax
}

func (l Level) String() string {
	if lmin < l && l < lmax {
		return levels[l]
	}

	return "UNKNOWN"
}

// Resolves level by name, returns lmin without definition by default
func ResolveLevelByName(name string) Level {
	switch strings.ToUpper(name) {
	case "DEBUG":
		return Ldebug

	case "INFO":
		return Linfo

	case "WARN":
		return Lwarn

	case "ERROR":
		return Lerror

	case "FATAL":
		return Lfatal

	case "PANIC":
		return Lpanic

	case "STACK":
		return Ltrace

	}

	return lmin
}
