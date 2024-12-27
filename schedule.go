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

// Schedule represents the runs occurring at a GDQ event.
type Schedule struct {
	Runs     []*Run
	byRunner map[string][]*Run
	byHost   map[string][]*Run
	l        sync.RWMutex
}

// NewSchedule returns an empty Schedule.
func NewSchedule() *Schedule {
	return &Schedule{
		Runs:     []*Run{},
		byRunner: map[string][]*Run{},
		byHost:   map[string][]*Run{},
	}
}

// NewScheduleFrom returns a scheduled filled with the runs.
func NewScheduleFrom(runs []*Run) *Schedule {
	if len(runs) == 0 {
		return nil
	}

	s := &Schedule{
		Runs:     runs,
		byRunner: map[string][]*Run{},
		byHost:   map[string][]*Run{},
	}

	s.calc()
	return s
}

// calc computes the byHost and byRunner lookup maps.
func (s *Schedule) calc() {
	s.l.Lock()
	defer s.l.Unlock()

	for _, run := range s.Runs {
		for _, talent := range run.Runners {
			name := normalised(talent.Name)
			if rev, ok := s.byRunner[name]; ok {
				s.byRunner[name] = append(rev, run)
			} else {
				s.byRunner[name] = []*Run{run}
			}
		}
		for _, talent := range run.Hosts {
			name := normalised(talent.Name)
			if hev, ok := s.byHost[name]; ok {
				s.byHost[name] = append(hev, run)
			} else {
				s.byHost[name] = []*Run{run}
			}
		}
	}
}

// ForRunner returns a new schedule with runs only matching this runner.
//
// The runner's name is matched using a substring match. This means that if you
// call somtething like schedule.ForRunner("b") you can get a schedule with runs
// for multiple runners.
//
// You'll get a nil schedule if no run matched the runner.
//
// The match is case insensitive.
func (s *Schedule) ForRunner(name string) *Schedule {
	return s.forEntity("runner", name)
}

// ForHost returns a new schedule with runs only matching this host.
//
// The host's name is matched using a substring match. This means that if you
// call somtething like schedule.ForHost("b") you can get a schedule with runs
// for multiple hosts.
//
// You'll get a nil schedule if no run matched the host.
//
// The match is case insensitive.
func (s *Schedule) ForHost(name string) *Schedule {
	return s.forEntity("host", name)
}

// ForTitle returns a new schedule with runs only matching the title.
//
// The title is matched using a substring match. This means that if you call
// somtething like schedule.ForTitle("b") you can get a schedule with multiple
// runs.
//
// You'll get a nil schedule if no run matched the title.
//
// The match is case insensitive.
func (s *Schedule) ForTitle(title string) *Schedule {
	return s.forEntity("title", title)
}

func (s *Schedule) forEntity(kind string, match string) *Schedule {
	if strings.TrimSpace(match) == "" {
		return nil
	}

	match = normalised(match)
	matched := make([]*Run, 0, 8)

	s.l.RLock()
	switch kind {
	case "title":
		for _, run := range s.Runs {
			if strings.Contains(normalised(run.Title), match) {
				matched = append(matched, run)
			}
		}
	case "host":
		for h, rs := range s.byHost {
			if strings.Contains(h, match) {
				matched = append(matched, rs...)
			}
		}
	case "runner":
		for h, rs := range s.byRunner {
			if strings.Contains(h, match) {
				matched = append(matched, rs...)
			}
		}
	default:
		panic(fmt.Sprintf("unsupported kind: %s in forEntity call", kind))
	}
	s.l.RUnlock()

	if len(matched) == 0 {
		return nil
	}

	return NewScheduleFrom(matched)
}

// NextRun returns the next run in the [Schedule].
//
// It returns the first run after t.
func (s *Schedule) NextRun(t time.Time) *Run {
	s.l.RLock()
	defer s.l.RUnlock()
	for _, run := range s.Runs {
		if run.Start.After(t) {
			return run
		}
	}
	return nil
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
