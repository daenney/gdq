package gdq

import (
	"time"
)

type runResp struct {
	ID     uint `json:"pk"`
	Fields struct {
		Name     string    `json:"display_name"`
		Start    time.Time `json:"starttime"`
		Estimate Duration  `json:"run_time"`
		Category string    `json:"category"`
		Console  string    `json:"console"`
		Year     *uint     `json:"release_year"`
		Runners  []uint    `json:"runners"`
	} `json:"fields"`
}

func (r runResp) toRun() Run {
	return Run{
		Title:    r.Fields.Name,
		Start:    r.Fields.Start,
		Estimate: r.Fields.Estimate,
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
