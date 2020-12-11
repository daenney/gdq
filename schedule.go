package gdq

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/secure/precis"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/anaskhan96/soup"
)

const scheduleURL = "https://gamesdonequick.com/schedule"

// Schedule represents the events occurring at a GDQ
type Schedule struct {
	Events   []*Event
	byRunner map[string][]*Event
	byHost   map[string][]*Event
	l        sync.RWMutex
}

// NewSchedule returns an empty Schedule
func NewSchedule() *Schedule {
	return &Schedule{
		Events:   []*Event{},
		byRunner: map[string][]*Event{},
		byHost:   map[string][]*Event{},
	}
}

// NewScheduleFrom returns a scheduled filled with the events
func NewScheduleFrom(events []*Event) *Schedule {
	if events == nil || len(events) == 0 {
		return NewSchedule()
	}

	s := &Schedule{
		Events:   make([]*Event, 0, len(events)),
		byRunner: map[string][]*Event{},
		byHost:   map[string][]*Event{},
	}
	s.load(events)
	return s
}

// load a series of events in the Schedule
//
// Call this method when wanting to add events to a schedule to ensure that
// the byRunner and byHost maps get updated. This permits the filter functions
// like ForHost and ForRunner to work
func (s *Schedule) load(events []*Event) {
	s.l.Lock()
	defer s.l.Unlock()
	for _, event := range events {
		s.Events = append(s.Events, event)
		for _, runner := range event.Runners {
			rev, ok := s.byRunner[runner]
			if ok {
				s.byRunner[runner] = append(rev, event)
			} else {
				s.byRunner[runner] = []*Event{event}
			}
		}
		for _, host := range event.Hosts {
			hev, ok := s.byHost[host]
			if ok {
				s.byHost[host] = append(hev, event)
			} else {
				s.byHost[host] = []*Event{event}
			}
		}
	}
}

// ForRunner returns a new schedule with events only matching this runner
//
// The runner's name is matched using a string submatch. This means that if you
// call somtething like schedule.ForRunner("b") you can get a schedule with events
// for multiple runners.
//
// The match is case insensitive.
func (s *Schedule) ForRunner(name string) *Schedule {
	return s.forEntity("runner", name)
}

// ForHost returns a new schedule with events only matching this host
//
// The host's name is matched using a string submatch. This means that if you
// call somtething like schedule.ForHust("b") you can get a schedule with events
// for multiple hosts.
//
// The match is case insensitive.
func (s *Schedule) ForHost(name string) *Schedule {
	return s.forEntity("host", name)
}

// ForTitle returns a new schedule with events only matching this event title
//
// The title is matched using a string submatch. This means that if you call
// somtething like schedule.ForTitle("b") you can get a schedule with multiple
// events.
//
// The match is case insensitive.
func (s *Schedule) ForTitle(title string) *Schedule {
	if strings.TrimSpace(title) == "" {
		return NewSchedule()
	}

	s.l.RLock()
	defer s.l.RUnlock()

	evs := []*Event{}
	for _, e := range s.Events {
		if strings.Contains(normalised(e.Title), normalised(title)) {
			evs = append(evs, e)
		}
	}

	ns := NewScheduleFrom(evs)
	return ns
}

func (s *Schedule) forEntity(kind string, match string) *Schedule {
	ns := NewSchedule()
	if strings.TrimSpace(match) == "" {
		return ns
	}

	var ev map[string][]*Event

	switch kind {
	case "host":
		ev = s.byHost
	case "runner":
		ev = s.byRunner
	default:
		panic(fmt.Sprintf("unsupported kind: %s in forEntity call", kind))
	}

	s.l.RLock()
	defer s.l.RUnlock()

	for h := range ev {
		if strings.Contains(normalised(h), normalised(match)) {
			ns.load(ev[h])
		}
	}

	return ns
}

// NextEvent returns the next/upcoming event in the schedule
func (s *Schedule) NextEvent() *Event {
	now := time.Now().UTC()
	var ev *Event

	s.l.RLock()
	defer s.l.RUnlock()
	for _, event := range s.Events {
		if event.Start.After(now) {
			ev = event
			break
		}
	}
	return ev
}

// GetSchedule returns the Schedule for a GDQ edition
//
// A client has to be passed in. Please make sure to configure your client
// correctly, so not http.DefaultClient. Be nice to server admins and make
// sure your client sets a User-Agent header.
func GetSchedule(id Edition, client *http.Client) (*Schedule, error) {
	if client == nil {
		return nil, fmt.Errorf("missing *http.Client")
	}

	resp, err := soup.GetWithClient(fmt.Sprintf("%s/%d", scheduleURL, id), client)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %w", err)
	}

	doc := soup.HTMLParse(resp)
	if doc.Error != nil {
		return nil, ErrInvalidSchedule
	}
	table := doc.Find("table", "id", "runTable")
	if table.Error != nil {
		return nil, ErrMissingSchedule
	}
	body := table.Find("tbody")
	if body.Error != nil {
		return nil, ErrMissingSchedule
	}

	rows := body.FindAll("tr")
	if len(rows) < 2 {
		return nil, ErrMissingSchedule
	}

	if len(rows)%2 != 0 {
		return nil, ErrInvalidSchedule
	}

	events := []*Event{}
	for i := 0; i < len(rows); i += 2 {
		event, err := eventFromHTML(rows[i], rows[i+1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse rows %s and %s as an event: %w", rows[i].HTML(), rows[i+1].HTML(), err)
		}
		events = append(events, event)
	}

	schedule := NewScheduleFrom(events)

	return schedule, nil
}

// normalised transforms a string to a variant that has punctuation and
// diacritics removed, and is mapped to lower case
func normalised(s string) string {
	s = runes.Remove(runes.In(unicode.Punct)).String(s)
	filter := precis.NewIdentifier(
		precis.LowerCase(),
		precis.AdditionalMapping(func() transform.Transformer {
			return transform.Chain(
				norm.NFD,
				runes.Remove(runes.In(unicode.Mn)))
		}),
		precis.Norm(norm.NFC),
	)

	normalised := []string{}
	for _, p := range strings.Fields(s) {
		res, _ := filter.String(p)
		normalised = append(normalised, res)
	}

	return strings.Join(normalised, " ")
}
