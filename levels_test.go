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

func Test_Level_ResolveLevelByName(t *testing.T) {
	assertion := assert.New(t)

	for level, name := range levels {
		assertion.Equal(level, ResolveLevelByName(name))
	}

	assertion.Equal(lmin, ResolveLevelByName("UNKNOWN"))
}
