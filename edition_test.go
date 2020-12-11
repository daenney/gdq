package gdq

import (
	"testing"
)

func TestEdition(t *testing.T) {
	t.Run("latest", func(t *testing.T) {
		assertEqual(t, Edition(0).String(), "latest")
	})
	t.Run("unknown", func(t *testing.T) {
		assertEqual(t, Edition(^uint(0)).String(), "unknown edition: 18446744073709551615")
	})
}

func TestGetEdition(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		_, ok := GetEdition("")
		assertEqual(t, ok, false)
	})
	t.Run("unknown", func(t *testing.T) {
		_, ok := GetEdition("AGDQ0")
		assertEqual(t, ok, false)
	})
	t.Run("known", func(t *testing.T) {
		e, ok := GetEdition("AgdQ2016")
		assertEqual(t, ok, true)
		assertEqual(t, e, AGDQ2016)
	})
}
