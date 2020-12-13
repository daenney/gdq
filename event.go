package gdq

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

// Duration represents a interval of time
type Duration struct {
	time.Duration
}

// MarshalJSON marshals a Duration to JSON
func (d Duration) MarshalJSON() (b []byte, err error) {
	return json.Marshal(d.Duration)
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

// Event represents a single event at a GDQ
type Event struct {
	Start    time.Time `json:"start"`
	Setup    Duration  `json:"setup"`
	Estimate Duration  `json:"estimate"`
	Runners  []string  `json:"runners"`
	Hosts    []string  `json:"hosts"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Platform string    `json:"platform"`
}

func eventFromHTML(rows ...soup.Root) (*Event, error) {
	tr1 := rows[0].FindAll("td")
	tr2 := rows[1].FindAll("td")

	if len(tr1) < 4 || len(tr2) < 3 {
		return nil, ErrUnexpectedData
	}

	e := &Event{}

	e.Start = toDateTime(tr1[0].Text())
	e.Title = strings.TrimSpace(tr1[1].Text())
	e.Runners = nicksToSlice(tr1[2].Text())
	e.Setup = toDuration(tr1[3].Text())

	e.Estimate = toDuration(tr2[0].Text())
	catPlat := strings.Split(tr2[1].Text(), "—") // Ceci n'est pas un -
	category := "unknown"
	platform := "unknown"
	if len(catPlat) == 2 {
		if c := strings.TrimSpace(catPlat[0]); c != "" {
			category = c
		}
		if p := strings.TrimSpace(catPlat[1]); p != "" {
			platform = p
		}
	}
	e.Category = category
	e.Platform = platform
	e.Hosts = nicksToSlice(tr2[2].Text())

	return e, nil
}

func nicksToSlice(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	data := strings.Split(input, ",")
	res := make([]string, len(data))
	for i, d := range data {
		res[i] = strings.TrimSpace(d)
	}

	return res
}

func toDateTime(input string) time.Time {
	res, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return time.Time{}
	}
	return res
}

func toDuration(input string) Duration {
	if strings.TrimSpace(input) == "" {
		return Duration{0}
	}
	elems := strings.Split(strings.TrimSpace(input), ":")
	if len(elems) != 3 {
		return Duration{0}
	}
	entry := fmt.Sprintf("%sh%sm%ss", strings.TrimSpace(elems[0]), strings.TrimSpace(elems[1]), strings.TrimSpace(elems[2]))
	res, err := time.ParseDuration(entry)
	if err != nil {
		return Duration{0}
	}
	return Duration{res}
}
