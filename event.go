package gdq

import (
	"fmt"
	"strings"
	"time"
)

type eventResp struct {
	ID             uint      `json:"id"`
	Short          string    `json:"short"`
	Name           string    `json:"name"`
	StartTime      time.Time `json:"datetime"`
	DonationAmount float64   `json:"amount"`
	DonationCount  uint64    `json:"donation_count"`
}

func (e eventResp) toEvent() *Event {
	return &Event{
		ID:    e.ID,
		Short: e.Short,
		Name:  e.Name,
		Year:  e.StartTime.Year(),
		Donations: Donation{
			Amount: e.DonationAmount,
			Count:  e.DonationCount,
		},
	}
}

type eventsResp struct {
	Results []eventResp `json:"results"`
}

func (e eventsResp) toEvents() []*Event {
	evs := make([]*Event, 0, len(e.Results))
	for _, r := range e.Results {
		evs = append(evs, r.toEvent())
	}
	return evs
}

// Event is the schedule ID of a GDQ event
type Event struct {
	ID    uint   `json:"id"`
	Short string `json:"short"`
	Name  string `json:"name"`
	Year  int    `json:"year"`

	Donations Donation `json:"donations"`
}

func (e Event) String() string {
	return fmt.Sprintf("%s (%d)", e.Name, e.Year)
}

// GetEventByName tries to find an event matching the input
func GetEventByName(input string) (ev *Event, found bool) {
	e, ok := eventsByName[strings.ToLower(input)]
	return &e, ok
}

// GetEventByID fetches the event by ID
func GetEventByID(id uint) (ev *Event, found bool) {
	e, ok := eventsByID[id]
	return &e, ok
}
