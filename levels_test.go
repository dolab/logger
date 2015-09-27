package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Level_String(t *testing.T) {
	assertion := assert.New(t)

	for level, s := range levels {
		assertion.Equal(s, level.String())
	}

	assertion.Equal("UNKNOWN", lmin.String())
	assertion.Equal("UNKNOWN", lmax.String())
}
