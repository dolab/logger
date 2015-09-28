package logger

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/dolab/colorize"
)

var (
	// default flag of source path format
	flag = log.Ldate | log.Lmicroseconds | log.Llongfile

	// Logger brushes
	brushes = map[Level]colorize.Colorize{
		Ldebug: colorize.New("gray"),
		Linfo:  colorize.New("cyan"),
		Lwarn:  colorize.New("yellow"),
		Lerror: colorize.New("magenta"),
		Lfatal: colorize.New("red"),
		Lpanic: colorize.New("red"),
		Ltrace: colorize.New("green"),
	}
)

type Logger struct {
	mux sync.Mutex
	out io.Writer
	buf []byte

	level Level
	tags  []string
	flag  int
	skip  int
	color bool
}

// Create a logger with the requested output. (default to stderr)
// Available output are [stdout|stderr|null|nil|path/to/file]
func New(output string) (*Logger, error) {
	colorful := (runtime.GOOS != "windows")

	switch output {
	case "stdout":
		return &Logger{
			out:   os.Stdout,
			flag:  flag,
			skip:  2,
			color: colorful,
		}, nil

	case "stderr":
		return &Logger{
			out:   os.Stderr,
			flag:  flag,
			skip:  2,
			color: colorful,
		}, nil

	default:
		if output == "null" || output == "nil" {
			output = os.DevNull
		}

		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("Failed to open log file %s: %v", output, err)
		}

		return &Logger{
			out:   file,
			flag:  flag,
			skip:  2,
			color: false,
		}, nil
	}

	return nil, ErrOutput
}

// New allocates a new Logger with given tags
func (l *Logger) New(tags ...string) *Logger {
	tmp := *l
	tmp.mux.Lock()
	tmp.buf = tmp.buf[:0]
	tmp.tags = tags
	tmp.mux.Unlock()

	return &tmp
}

// SetLevel sets min level of output
func (l *Logger) SetLevel(level Level) error {
	if !level.IsValid() {
		return ErrLevel
	}

	l.mux.Lock()
	l.level = level
	l.mux.Unlock()

	return nil
}

// SetLevelByName sets min level of output by name,
// available values are [debug|info|warn|error|fatal|panic|stack].
// It returns ErrLevel for invalid name.
func (l *Logger) SetLevelByName(name string) error {
	level := ResolveLevelByName(name)
	if !level.IsValid() {
		return ErrLevel
	}

	l.mux.Lock()
	l.level = level
	l.mux.Unlock()

	return nil
}

// SetTags sets tags of all logs, it'll replace previous definition.
func (l *Logger) SetTags(tags ...string) {
	l.mux.Lock()
	l.tags = tags
	l.mux.Unlock()
}

// AddTags adds new tags to all logs, duplicated tags will be ignored.
func (l *Logger) AddTags(tags ...string) {
	l.mux.Lock()
	newTags := []string{}
	for _, tag := range tags {
		found := false
		for _, existedTag := range l.tags {
			if tag == existedTag {
				found = true
				break
			}
		}

		if !found {
			newTags = append(newTags, tag)
		}
	}

	if len(newTags) > 0 {
		l.tags = append(l.tags, newTags...)
	}
	l.mux.Unlock()
}

// SetFlag changes flag of source file path format
func (l *Logger) SetFlag(flag int) {
	l.mux.Lock()
	l.flag = flag
	l.mux.Unlock()
}

// SetSkip changes the PC
func (l *Logger) SetSkip(depth int) {
	l.mux.Lock()
	l.skip = depth
	l.mux.Unlock()
}

// SetColor sets whether output logs with colorful
func (l *Logger) SetColor(colorful bool) {
	l.mux.Lock()
	l.color = colorful
	l.mux.Unlock()
}

// SetOutput sets output of Logger
func (l *Logger) SetOutput(w io.Writer) {
	l.mux.Lock()
	l.out = w
	l.mux.Unlock()
}

// Output writes the output for a logging event.
// The string s contains the text to print after the tags specified
// by the flags of the Logger.
// A newline is appended if the last character of s is not already a newline.
func (l *Logger) Output(level Level, s string) error {
	if !level.IsValid() {
		return ErrLevel
	}

	l.mux.Lock()
	defer l.mux.Unlock()

	var (
		file string
		line int
		ok   bool
	)

	if l.flag&(log.Lshortfile|log.Llongfile) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mux.Unlock()

		_, file, line, ok = runtime.Caller(l.skip)
		if !ok {
			file = "???"
			line = 0
		}

		l.mux.Lock()
	}

	l.buf = l.buf[:0]
	l.formatHeader(level, file, line, &l.buf)
	l.buf = append(l.buf, s...)

	// append newline if needs
	if len(s) > 0 && s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	if l.color {
		l.buf = []byte(brushes[level].Paint(string(l.buf)))
	}

	_, err := l.out.Write(l.buf)
	return err
}

// Debug calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Debug(v ...interface{}) {
	if l.level > Ldebug {
		return
	}

	l.Output(Ldebug, fmt.Sprint(v...))
}

// Debugf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level > Ldebug {
		return
	}

	l.Output(Ldebug, fmt.Sprintf(format, v...))
}

// Info calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	if l.level > Linfo {
		return
	}

	l.Output(Linfo, fmt.Sprint(v...))
}

// Infof calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level > Linfo {
		return
	}

	l.Output(Linfo, fmt.Sprintf(format, v...))
}

// Warn calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Warn(v ...interface{}) {
	if l.level > Lwarn {
		return
	}

	l.Output(Lwarn, fmt.Sprint(v...))
}

// Warnf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level > Lwarn {
		return
	}

	l.Output(Lwarn, fmt.Sprintf(format, v...))
}

// Error calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	if l.level > Lerror {
		return
	}

	l.Output(Lerror, fmt.Sprint(v...))
}

// Errorf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level > Lerror {
		return
	}

	l.Output(Lerror, fmt.Sprintf(format, v...))
}

// Fatal calls l.Output to print to the logger and exit process with sign 1.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Fatal(v ...interface{}) {
	if l.level > Lfatal {
		return
	}

	l.Output(Lfatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf calls l.Output to print to the logger and exit process with sign 1.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.level > Lfatal {
		return
	}

	l.Output(Lfatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Panic calls l.Output to print to the logger and panic process.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Panic(v ...interface{}) {
	if l.level > Lpanic {
		return
	}

	s := fmt.Sprint(v...)
	l.Output(Lpanic, s)
	panic(s)
}

// Panicf calls l.Output to print to the logger and panic process.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Panicf(format string, v ...interface{}) {
	if l.level > Lpanic {
		return
	}

	s := fmt.Sprintf(format, v...)
	l.Output(Lpanic, s)
	panic(s)
}

// Trace calls l.Output to print to the logger and output process stacks,
// and exit process with sign 1 at last.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Trace(v ...interface{}) {
	l.Output(Ltrace, fmt.Sprint(v...))

	// process stacks
	brush := brushes[Ltrace]
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)

	scanner := bufio.NewScanner(bytes.NewReader(buf[:n]))
	for scanner.Scan() {
		line := scanner.Text()

		if l.color {
			line = brush.Paint(line)
		}

		l.buf = append(l.buf, line...)
		l.buf = append(l.buf, '\n')
	}

	l.out.Write([]byte(string(l.buf)))
	os.Exit(1)
}

// Modified from src/log/log.go
func (l *Logger) formatHeader(level Level, file string, line int, buf *[]byte) {
	t := time.Now()

	if l.flag&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		if l.flag&log.Ldate != 0 {
			year, month, day := t.Date()

			itoa(buf, year, 4)
			*buf = append(*buf, '/')

			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')

			itoa(buf, day, 2)
		}

		if l.flag&(log.Ltime|log.Lmicroseconds) != 0 {
			*buf = append(*buf, ' ')

			hour, min, sec := t.Clock()

			itoa(buf, hour, 2)
			*buf = append(*buf, ':')

			itoa(buf, min, 2)
			*buf = append(*buf, ':')

			itoa(buf, sec, 2)
			if l.flag&log.Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
		}

		*buf = append(*buf, " - "...)
	}

	*buf = append(*buf, '[')
	*buf = append(*buf, level.String()...)
	if len(l.tags) > 0 {
		*buf = append(*buf, ", "...)
		*buf = append(*buf, strings.Join(l.tags, ", ")...)
	}
	*buf = append(*buf, ']')
	*buf = append(*buf, " - "...)

	if l.flag&(log.Lshortfile|log.Llongfile) != 0 {
		short := file
		if l.flag&log.Lshortfile != 0 {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
		} else {
			for i := 0; i < len(file)-5; i++ {
				if file[i:i+5] == "/src/" {
					short = file[i+1:]
					break
				}
			}
		}

		*buf = append(*buf, short...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

// Cheap integer to fixed-width decimal ASCII.
// Give a negative width to avoid zero-padding.
// Knows the buffer has capacity.
func itoa(buf *[]byte, i int, wid int) {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		*buf = append(*buf, '0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}

	*buf = append(*buf, b[bp:]...)
}
