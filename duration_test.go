package gdq

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
)

func TestDurationMarshal(t *testing.T) {
	d := Duration{0}
	j, err := d.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(j), "0")

	dn := Duration{1*time.Hour + 1*time.Minute + 1*time.Second}
	jn, err := dn.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(jn), "3661000000000")
}

func TestDurationUnmarshall(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		d := &Duration{}
		err := d.UnmarshalJSON([]byte(`16000`))
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(16000), d.Duration)
	})
	t.Run("string h:m:s", func(t *testing.T) {
		d := &Duration{}
		err := d.UnmarshalJSON([]byte(`"00:10:00"`))
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(10*time.Minute), d.Duration)
	})
	t.Run("string 1:2", func(t *testing.T) {
		d := &Duration{}
		err := d.UnmarshalJSON([]byte(`"10:00"`))
		assert.Error(t, err)
	})
	t.Run("string a:b:c", func(t *testing.T) {
		d := &Duration{}
		err := d.UnmarshalJSON([]byte(`"a:b:c"`))
		assert.Error(t, err)
	})
	t.Run("string time.Duration", func(t *testing.T) {
		d := &Duration{}
		err := d.UnmarshalJSON([]byte(`"10m"`))
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(10*time.Minute), d.Duration)
	})
	t.Run("string other", func(t *testing.T) {
		d := &Duration{}
		err := d.UnmarshalJSON([]byte(`"banana"`))
		assert.Error(t, err)
	})
	t.Run("other type", func(t *testing.T) {
		d := &Duration{}
		err := d.UnmarshalJSON([]byte(`[]`))
		assert.Error(t, err)
	})
	t.Run("invalid JSON", func(t *testing.T) {
		d := &Duration{}
		err := d.UnmarshalJSON([]byte(`banana`))
		assert.Error(t, err)
	})
}

func TestDurationString(t *testing.T) {
	t.Run("no hours, no minutes", func(t *testing.T) {
		assert.Equal(t, Duration{0}.String(), "0 minutes")
	})
	t.Run("one hour, no minutes", func(t *testing.T) {
		assert.Equal(t, Duration{1 * time.Hour}.String(), "1 hour")
	})
	t.Run("one hour, one minute", func(t *testing.T) {
		assert.Equal(t, Duration{1*time.Hour + 1*time.Minute}.String(), "1 hour and 1 minute")
	})
	t.Run("one hour, two minutes", func(t *testing.T) {
		assert.Equal(t, Duration{1*time.Hour + 2*time.Minute}.String(), "1 hour and 2 minutes")
	})
	t.Run("two hours, two minutes", func(t *testing.T) {
		assert.Equal(t, Duration{2*time.Hour + 2*time.Minute}.String(), "2 hours and 2 minutes")
	})
	t.Run("one minute", func(t *testing.T) {
		assert.Equal(t, Duration{1 * time.Minute}.String(), "1 minute")
	})
	t.Run(" two minutes", func(t *testing.T) {
		assert.Equal(t, Duration{2 * time.Minute}.String(), "2 minutes")
	})
}
