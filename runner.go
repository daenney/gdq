package gdq

import "strings"

type runnerResp struct {
	ID     uint `json:"pk"`
	Fields struct {
		Handle  string `json:"public"`
		Stream  string `json:"stream"`
		Twitter string `json:"twitter"`
		YouTube string `json:"youtube"`
	} `json:"fields"`
}

func (r runnerResp) toRunner() Runner {
	return Runner{
		Handle: r.Fields.Handle,
		Social: social{
			Stream:  r.Fields.Stream,
			Twitter: r.Fields.Twitter,
			YouTube: r.Fields.YouTube,
		},
	}
}

type runnersResp []runnerResp

// Runner represents a person running a game
type Runner struct {
	Handle string `json:"handle"`
	Social social `json:"social"`
}

type social struct {
	Stream  string `json:"stream,omitempty"`
	Twitter string `json:"twitter,omitempty"`
	YouTube string `json:"youtube,omitempty"`
}

type Runners []Runner

func (r Runners) String() string {
	handles := []string{}

	for _, runner := range r {
		handles = append(handles, runner.Handle)
	}

	return strings.Join(handles, ", ")
}
