package gdq

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Duration represents a interval of time
type Duration struct {
	time.Duration
}

// MarshalJSON marshals a Duration to JSON
func (d Duration) MarshalJSON() (b []byte, err error) {
	return json.Marshal(d.Duration)
}

// UnmarshalJSON unmarshals a duration like thing from JSON
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		if strings.Contains(value, ":") {
			hms := strings.Split(value, ":")
			if len(hms) != 3 {
				return fmt.Errorf("invalid duration")
			}
			d.Duration, err = time.ParseDuration(fmt.Sprintf("%sh%sm%ss", hms[0], hms[1], hms[2]))
			if err != nil {
				return err
			}
			return nil
		}
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("invalid duration")
	}
}

func (d Duration) Add(nd Duration) Duration {
	return Duration{time.Duration(d.Nanoseconds() + nd.Nanoseconds())}
}

func (d Duration) String() string {
	dur := d.Round(time.Minute)
	h := dur / time.Hour
	m := (dur - (h * time.Hour)) / time.Minute
	if h == 0 {
		if m == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", m)
	}
	b := strings.Builder{}
	if h == 1 {
		b.WriteString("1 hour")
	} else {
		b.WriteString(fmt.Sprintf("%d hours", h))
	}

	if m == 0 {
		return b.String()
	}

	b.WriteString(" and ")
	if m == 1 {
		b.WriteString("1 minute")
	} else {
		b.WriteString(fmt.Sprintf("%d minutes", m))
	}
	return b.String()
}
