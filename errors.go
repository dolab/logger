package logger

import (
	"errors"
)

var (
	ErrOutput = errors.New("Unsupported output")
	ErrLevel  = errors.New("Invalid level")
)
