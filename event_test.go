package gdq

import (
	"math"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestGetEvent(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		_, ok := GetEventByName("")
		assert.Equal(t, false, ok)
	})
	t.Run("unknown", func(t *testing.T) {
		_, ok := GetEventByName("AGDQ0")
		assert.Equal(t, false, ok)
	})
	t.Run("known", func(t *testing.T) {
		e, ok := GetEventByName("AgdQ2016")
		assert.True(t, ok)
		assert.Equal(t, &AGDQ2016, e)
	})
}

func TestGetEventByID(t *testing.T) {
	t.Run("unknown", func(t *testing.T) {
		_, ok := GetEventByID(math.MaxInt64)
		assert.Equal(t, false, ok)
	})
	t.Run("known", func(t *testing.T) {
		ev, ok := GetEventByID(17)
		assert.True(t, ok)
		assert.Equal(t, &AGDQ2016, ev)
	})
}

func TestEventString(t *testing.T) {
	assert.Equal(t, "Awesome Games Done Quick (2016)", AGDQ2016.String())
}
