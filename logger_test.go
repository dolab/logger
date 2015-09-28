package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
