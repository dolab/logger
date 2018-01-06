package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/golib/assert"
)

func Test_Logger_New(t *testing.T) {
	logger, err := New("nil")
	assert.Nil(t, err)
	assert.Implements(t, (*io.Writer)(nil), logger)
}

func Test_Logger_NewWithStdout(t *testing.T) {
	logger, err := New("stdout")
	assert.Nil(t, err)
	assert.Equal(t, os.Stdout, logger.out)
}

func Test_Logger_NewWithStderr(t *testing.T) {
	logger, err := New("stderr")
	assert.Nil(t, err)
	assert.Equal(t, os.Stderr, logger.out)
}

func Test_Logger_NewWithDevNull(t *testing.T) {
	logger, err := New("null")
	assert.Nil(t, err)
	assert.IsType(t, (*os.File)(nil), logger.out)

	logger, err = New("nil")
	assert.Nil(t, err)
	assert.IsType(t, (*os.File)(nil), logger.out)
}

func Test_Logger_NewWithFile(t *testing.T) {
	filename := "logger-testing"

	logger, err := New(filename)
	assert.Nil(t, err)
	assert.IsType(t, (*os.File)(nil), logger.out)

	os.Remove(filename)
}

func Test_Logger_NewWithTags(t *testing.T) {
	// mock os.Stdout
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	tags := []string{"testing", "logger"}
	s := "output testing"
	expected := "[DEBUG, testing, logger]"

	logger, _ := New("stdout")
	logger.SetSkip(1)

	taggedLogger := logger.New(tags...)
	taggedLogger.Debug(s)

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	assert.Nil(t, err)
	assert.Contains(t, string(buf[:n]), expected)

	os.Stdout = stdout
}

func Test_Logger_SetTags(t *testing.T) {
	// mock os.Stdout
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	tags := []string{"testing", "logger"}
	s := "output testing"
	expected := "[DEBUG, testing, logger]"

	logger, _ := New("stdout")
	logger.SetSkip(1)
	logger.SetTags(tags...)

	logger.Debug(s)

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	assert.Nil(t, err)
	assert.Contains(t, string(buf[:n]), expected)

	os.Stdout = stdout
}

func Test_Logger_AddTags(t *testing.T) {
	// mock os.Stdout
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	tags := []string{"testing", "logger"}
	newTags := []string{"new_testing", "logger"}
	s := "output testing"
	expected := "[DEBUG, testing, logger, new_testing]"

	logger, _ := New("stdout")
	logger.SetSkip(1)
	logger.SetTags(tags...)
	logger.AddTags(newTags...)

	logger.Debug(s)

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	assert.Nil(t, err)
	assert.Contains(t, string(buf[:n]), expected)

	os.Stdout = stdout
}

func Test_Logger_Output(t *testing.T) {
	// mock os.Stdout
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	s := "output testing"
	expected := s + "\n"
	testCases := map[Level]string{
		Ldebug: "[DEBUG]",
		Linfo:  "[INFO]",
		Lwarn:  "[WARN]",
		Lerror: "[ERROR]",
		Lfatal: "[FATAL]",
		Lpanic: "[PANIC]",
		Ltrace: "[Stack]",
	}

	logger, _ := New("stdout")
	logger.SetSkip(1)

	for level, tag := range testCases {
		logger.Output(level, s)

		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		assert.Nil(t, err)
		assert.Contains(t, string(buf[:n]), tag)
		assert.Contains(t, string(buf[:n]), expected)
	}

	os.Stdout = stdout

	logger.Output(Linfo, "Hello, logger!")
	logger.Output(Linfo, "Hello, logger!")
}

func Benchmark_Logger_Output(b *testing.B) {
	logger, _ := New("stderr")
	logger.SetSkip(1)
	logger.SetColor(false)
	logger.Output(Linfo, "Hello, logger!")

	logger.SetOutput(ioutil.Discard)

	for i := 0; i < b.N; i++ {
		logger.Output(Linfo, "Hello, logger!")
	}
}

func Benchmark_Logger_Stdlib(b *testing.B) {
	log.SetFlags(log.Llongfile | log.Ltime)
	log.SetPrefix("[INFO]")
	log.Println("Hello, logger!")

	log.SetOutput(ioutil.Discard)

	for i := 0; i < b.N; i++ {
		log.Println("Hello, logger!")
	}
}

func Test_Logger_Lock(t *testing.T) {
	logger, _ := New("stdout")
	logger.SetSkip(1)

	var (
		wg sync.WaitGroup

		routines = 10
	)
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func(routine int) {
			defer wg.Done()

			log := logger.New()
			assert.NotContains(t, strings.Join(log.Tags(), ", "), "new logger")

			log.SetTags("new logger")
			assert.Contains(t, strings.Join(log.Tags(), ", "), "new logger")
		}(i)
	}

	wg.Wait()
}

func Test_Logger_Racy(t *testing.T) {
	logger, _ := New("stdout")
	logger.SetSkip(1)

	var (
		wg sync.WaitGroup

		routines = 10
	)
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func(routine int) {
			defer wg.Done()

			logger.Infof("[OK] routine@#%d", routine)
		}(i)
	}

	wg.Wait()
}
