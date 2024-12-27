package gdq

import (
	"time"
)

type runResp struct {
	Results []struct {
		ID           uint      `json:"id"`
		Name         string    `json:"name"`
		Category     string    `json:"category"`
		Console      string    `json:"console"`
		Runners      []Talent  `json:"runners"`
		Hosts        []Talent  `json:"hosts"`
		Commentators []Talent  `json:"commentators"`
		Starttime    time.Time `json:"starttime"`
		Endtime      time.Time `json:"endtime"`
		RunTime      Duration  `json:"run_time"`
		SetupTime    Duration  `json:"setup_time"`
	} `json:"results"`
}

func (r runResp) toRuns() []*Run {
	lrun := len(r.Results)
	if lrun == 0 {
		return nil
	}

	runs := make([]*Run, 0, lrun)
	for _, r := range r.Results {
		if r.RunTime.Milliseconds() == 0 {
			// A lot of older events have their Estimate always set to 0, and the
			// same for their setup time. When we run into that, subtract the
			// start time from the end time. It's not 100% accurate since it doesn't
			// account for the setup time, but it's better than just showing 0 everywhere
			r.RunTime = Duration{r.Endtime.Sub(r.Starttime)}
		}
		runs = append(runs, &Run{
			Title:        r.Name,
			Start:        r.Starttime,
			Estimate:     r.RunTime.Add(r.SetupTime),
			Category:     r.Category,
			Platform:     r.Console,
			Hosts:        r.Hosts,
			Runners:      r.Runners,
			Commentators: r.Commentators,
		})
	}
	return runs
}

// Run represents a single event at a GDQ
type Run struct {
	Title        string    `json:"title"`
	Start        time.Time `json:"start"`
	Estimate     Duration  `json:"estimate"`
	Runners      []Talent  `json:"runners"`
	Hosts        []Talent  `json:"hosts"`
	Commentators []Talent  `json:"commentators"`
	Category     string    `json:"category"`
	Platform     string    `json:"platform"`
}
