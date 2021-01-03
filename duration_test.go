package gdq

import (
	"testing"
	"time"
)

func TestDurationMarshal(t *testing.T) {
	d := Duration{0}
	j, err := d.MarshalJSON()
	assertEqual(t, err, nil)
	assertEqual(t, string(j), "0")

	dn := Duration{1*time.Hour + 1*time.Minute + 1*time.Second}
	jn, err := dn.MarshalJSON()
	assertEqual(t, err, nil)
	assertEqual(t, string(jn), "3661000000000")
}

func TestDurationString(t *testing.T) {
	t.Run("no hours, no minutes", func(t *testing.T) {
		assertEqual(t, Duration{0}.String(), "0 minutes")
	})
	t.Run("one hour, no minutes", func(t *testing.T) {
		assertEqual(t, Duration{1 * time.Hour}.String(), "1 hour")
	})
	t.Run("one hour, one minute", func(t *testing.T) {
		assertEqual(t, Duration{1*time.Hour + 1*time.Minute}.String(), "1 hour and 1 minute")
	})
	t.Run("one hour, two minutes", func(t *testing.T) {
		assertEqual(t, Duration{1*time.Hour + 2*time.Minute}.String(), "1 hour and 2 minutes")
	})
	t.Run("two hours, two minutes", func(t *testing.T) {
		assertEqual(t, Duration{2*time.Hour + 2*time.Minute}.String(), "2 hours and 2 minutes")
	})
	t.Run("one minute", func(t *testing.T) {
		assertEqual(t, Duration{1 * time.Minute}.String(), "1 minute")
	})
	t.Run(" two minutes", func(t *testing.T) {
		assertEqual(t, Duration{2 * time.Minute}.String(), "2 minutes")
	})
}
