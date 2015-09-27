package logger

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

func (l Level) String() string {
	if lmin < l && l < lmax {
		return levels[l]
	}

	return "UNKNOWN"
}
