package gdq

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

// Event represents a single event at a GDQ
type Event struct {
	Start    time.Time
	Setup    time.Duration
	Estimate time.Duration
	Runners  []string
	Hosts    []string
	Title    string
	Category string
	Platform string
}

// GetSchedule returns a list of GDQ events for a GDQ edition
func GetSchedule(id Edition, client *http.Client) ([]*Event, error) {
	if client == nil {
		client = newHTTPClient()
	}

	resp, err := soup.GetWithClient(fmt.Sprintf("https://gamesdonequick.com/schedule/%d", id), client)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %w", err)
	}

	doc := soup.HTMLParse(resp)
	rows := doc.Find("table", "id", "runTable").Find("tbody").FindAll("tr")
	if len(rows) < 2 {
		return nil, ErrMissingSchedule
	}

	if len(rows)%2 != 0 {
		return nil, ErrInvalidSchedule
	}

	events := []*Event{}
	for i := 0; i < len(rows); i += 2 {
		event, err := createEvent(rows[i], rows[i+1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse rows %s and %s as an event: %w", rows[i].HTML(), rows[i+1].HTML(), err)
		}
		events = append(events, event)
	}

	return events, nil
}

func createEvent(rows ...soup.Root) (*Event, error) {
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

func toDuration(input string) time.Duration {
	elems := strings.Split(input, ":")
	entry := fmt.Sprintf("%sh%sm%ss", strings.TrimSpace(elems[0]), strings.TrimSpace(elems[1]), strings.TrimSpace(elems[2]))
	res, err := time.ParseDuration(entry)
	if err != nil {
		return 0
	}
	return res
}
