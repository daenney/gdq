package gdq

import (
	"testing"
)

func TestGetEvent(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		_, ok := GetEventByName("")
		assertEqual(t, ok, false)
	})
	t.Run("unknown", func(t *testing.T) {
		_, ok := GetEventByName("AGDQ0")
		assertEqual(t, ok, false)
	})
	t.Run("known", func(t *testing.T) {
		e, ok := GetEventByName("AgdQ2016")
		assertEqual(t, ok, true)
		assertEqual(t, e, AGDQ2016)
	})
}

func TestEventString(t *testing.T) {
	assertEqual(t, "Awesome Games Done Quick (2016)", AGDQ2016.String())
}
