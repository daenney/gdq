package gdq

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/secure/precis"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Schedule represents the runs occurring at a GDQ event
type Schedule struct {
	Runs     []*Run
	byRunner map[string][]*Run
	byHost   map[string][]*Run
	l        sync.RWMutex
}

// NewSchedule returns an empty Schedule
func NewSchedule() *Schedule {
	return &Schedule{
		Runs:     []*Run{},
		byRunner: map[string][]*Run{},
		byHost:   map[string][]*Run{},
	}
}

// NewScheduleFrom returns a scheduled filled with the runs
func NewScheduleFrom(runs []*Run) *Schedule {
	if runs == nil || len(runs) == 0 {
		return NewSchedule()
	}

	s := &Schedule{
		Runs:     make([]*Run, 0, len(runs)),
		byRunner: map[string][]*Run{},
		byHost:   map[string][]*Run{},
	}
	s.load(runs)
	return s
}

// load a series of runs in the Schedule
//
// Call this method when wanting to add runs to a schedule to ensure that
// the byRunner and byHost maps get updated. This permits the filter functions
// like ForHost and ForRunner to work
func (s *Schedule) load(runs []*Run) {
	s.l.Lock()
	defer s.l.Unlock()
	for _, run := range runs {
		s.Runs = append(s.Runs, run)
		for _, runner := range run.Runners {
			rev, ok := s.byRunner[runner.Handle]
			if ok {
				s.byRunner[runner.Handle] = append(rev, run)
			} else {
				s.byRunner[runner.Handle] = []*Run{run}
			}
		}
		for _, host := range run.Hosts {
			hev, ok := s.byHost[host]
			if ok {
				s.byHost[host] = append(hev, run)
			} else {
				s.byHost[host] = []*Run{run}
			}
		}
	}
}

// ForRunner returns a new schedule with runs only matching this runner
//
// The runner's name is matched using a string submatch. This means that if you
// call somtething like schedule.ForRunner("b") you can get a schedule with runs
// for multiple runners.
//
// The match is case insensitive.
func (s *Schedule) ForRunner(name string) *Schedule {
	return s.forEntity("runner", name)
}

// ForHost returns a new schedule with runs only matching this host
//
// The host's name is matched using a string submatch. This means that if you
// call somtething like schedule.ForHust("b") you can get a schedule with runs
// for multiple hosts.
//
// The match is case insensitive.
func (s *Schedule) ForHost(name string) *Schedule {
	return s.forEntity("host", name)
}

// ForTitle returns a new schedule with runs only matching this runs title
//
// The title is matched using a string submatch. This means that if you call
// somtething like schedule.ForTitle("b") you can get a schedule with multiple
// runs.
//
// The match is case insensitive.
func (s *Schedule) ForTitle(title string) *Schedule {
	if strings.TrimSpace(title) == "" {
		return NewSchedule()
	}

	s.l.RLock()
	defer s.l.RUnlock()

	runs := []*Run{}
	for _, run := range s.Runs {
		if strings.Contains(normalised(run.Title), normalised(title)) {
			runs = append(runs, run)
		}
	}

	ns := NewScheduleFrom(runs)
	return ns
}

func (s *Schedule) forEntity(kind string, match string) *Schedule {
	ns := NewSchedule()
	if strings.TrimSpace(match) == "" {
		return ns
	}

	var runs map[string][]*Run

	switch kind {
	case "host":
		runs = s.byHost
	case "runner":
		runs = s.byRunner
	default:
		panic(fmt.Sprintf("unsupported kind: %s in forEntity call", kind))
	}

	s.l.RLock()
	defer s.l.RUnlock()

	for h := range runs {
		if strings.Contains(normalised(h), normalised(match)) {
			ns.load(runs[h])
		}
	}

	return ns
}

// NextRun returns the next/upcoming run in the schedule
func (s *Schedule) NextRun() *Run {
	now := time.Now().UTC()
	var runs *Run

	s.l.RLock()
	defer s.l.RUnlock()
	for _, run := range s.Runs {
		if run.Start.After(now) {
			runs = run
			break
		}
	}
	return runs
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
