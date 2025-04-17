package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/golib/assert"
)

func TestSlog(t *testing.T) {
	// mock os.Stdout
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger, _ := New("stdout")

	logger.LogAttrs(context.Background(), slog.LevelDebug, "hello world", slog.String("key", "value"))

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	assert.Nil(t, err)
	assert.Contains(t, string(buf[:n]), "key=value,")
	assert.Contains(t, string(buf[:n]), "msg=hello world")

	os.Stdout = stdout

	slog.Log(context.Background(), LevelDebug, "hello world", slog.String("key", "value"))
}
