package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/golib/assert"
)

func Test_Logger_NewTextLogger(t *testing.T) {
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

	logger.NewTextLogger().Str("key", "value").Bool("bool", true).Debug(s)

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	assert.Nil(t, err)
	assert.Contains(t, string(buf[:n]), expected)
	assert.Contains(t, string(buf[:n]), "key=value,")
	assert.Contains(t, string(buf[:n]), "bool=true,")
	assert.Contains(t, string(buf[:n]), "msg=output testing")

	os.Stdout = stdout
}

func Test_Logger_NewJsonLogger(t *testing.T) {
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

	logger.NewJsonLogger().Str("key", "value").Bool("bool", true).Debug(s)

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	assert.Nil(t, err)
	assert.Contains(t, string(buf[:n]), expected)
	assert.Contains(t, string(buf[:n]), `{"bool":true,"key":"value"}`)
	assert.Contains(t, string(buf[:n]), "output testing")
	assert.NotContains(t, string(buf[:n]), "msg=output testing")

	os.Stdout = stdout
}

func Test_Logger_Error(t *testing.T) {
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

	logger.NewJsonLogger().Str("key", "value").Err(fmt.Errorf("debugging"), false).Debug(s)

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	assert.Nil(t, err)
	assert.Contains(t, string(buf[:n]), expected)
	assert.Contains(t, string(buf[:n]), `{"error":"debugging","key":"value"}`)
	assert.Contains(t, string(buf[:n]), "output testing")
	assert.NotContains(t, string(buf[:n]), "msg=output testing")

	os.Stdout = stdout
}

func Test_Logger_Panic(t *testing.T) {
	logger, _ := New("stdout")
	assert.NotPanics(t, func() {
		logger.NewJsonLogger().Bool("key", true).Any("nil", nil).Err(nil, true).Error("panic")
	})
	assert.NotPanics(t, func() {
		logger.NewTextLogger().Bool("key", false).Any("nil", nil).Err(nil, true).Error("panic")
	})
}
