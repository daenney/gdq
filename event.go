package gdq

import (
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

type duration struct {
	time.Duration
}

func (d duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

// Event represents a single event at a GDQ
type Event struct {
	Start    time.Time `json:"start"`
	Setup    duration  `json:"setup"`
	Estimate duration  `json:"estimate"`
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
	e.Category = strings.TrimSpace(catPlat[0])
	e.Platform = strings.TrimSpace(catPlat[1])
	e.Hosts = nicksToSlice(tr2[2].Text())

	return e, nil
}

func nicksToSlice(input string) []string {
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

func toDuration(input string) duration {
	if strings.TrimSpace(input) == "" {
		return duration{0}
	}
	elems := strings.Split(input, ":")
	entry := fmt.Sprintf("%sh%sm%ss", strings.TrimSpace(elems[0]), strings.TrimSpace(elems[1]), strings.TrimSpace(elems[2]))
	res, err := time.ParseDuration(entry)
	if err != nil {
		return duration{0}
	}
	return duration{res}
}
