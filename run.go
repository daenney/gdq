package gdq

import (
	"time"
)

type runResp struct {
	ID     uint `json:"pk"`
	Fields struct {
		Name     string    `json:"display_name"`
		Start    time.Time `json:"starttime"`
		End      time.Time `json:"endtime"`
		Estimate Duration  `json:"run_time"`
		Category string    `json:"category"`
		Console  string    `json:"console"`
		Year     *uint     `json:"release_year"`
		Runners  []uint    `json:"runners"`
	} `json:"fields"`
}

func (r runResp) toRun() Run {
	est := r.Fields.Estimate
	// A lot of older events have their Estimate always set to 0, and the
	// same for their setup time. When we run into that, subtract the
	// start time from the end time. It's not 100% accurate since it doesn't
	// account for the setup time, but it's better than just showing 0 everywhere
	if est.Milliseconds() == 0 {
		est = Duration{r.Fields.End.Sub(r.Fields.Start)}
	}
	return Run{
		Title:    r.Fields.Name,
		Start:    r.Fields.Start,
		Estimate: est,
		Category: r.Fields.Category,
		Console:  r.Fields.Console,
		Hosts:    []string{},
		Runners:  []Runner{},
	}
}

type runsResp []runResp

// Run represents a single event at a GDQ
type Run struct {
	Title    string    `json:"title"`
	Start    time.Time `json:"start"`
	Estimate Duration  `json:"estimate"`
	Runners  Runners   `json:"runners"`
	Hosts    []string  `json:"hosts"`
	Category string    `json:"category"`
	Console  string    `json:"console"`
}
