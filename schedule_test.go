package gdq

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
)

var testRuns = []*Run{
	{
		Title:   "Game 1",
		Runners: []Talent{{Name: "amazing"}},
		Hosts:   []Talent{{Name: "wonderful"}},
		Start:   time.Now().Add(-10 * time.Minute).UTC(),
	},
	{
		Title:   "Game 2",
		Runners: []Talent{{Name: "amazing"}},
	},
	{
		Title: "Game 3",
		Hosts: []Talent{{Name: "wonderful"}},
	},
	{
		Title:   "Game 4",
		Runners: []Talent{{Name: "amazing"}},
		Hosts:   []Talent{{Name: "awesome"}},
		Start:   time.Now().Add(10 * time.Minute).UTC(),
	},
	{
		Title:   "Game 5",
		Runners: []Talent{{Name: "fantastic"}},
		Hosts:   []Talent{{Name: "wonderful"}},
		Start:   time.Now().Add(30 * time.Minute).UTC(),
	},
}

func TestNewScheduleFrom(t *testing.T) {
	t.Run("no runs", func(t *testing.T) {
		runs := []*Run{}
		s := NewScheduleFrom(runs)
		assert.Equal(t, nil, s)
	})
	t.Run("single empty run", func(t *testing.T) {
		runs := []*Run{{}}
		s := NewScheduleFrom(runs)
		assert.Equal(t, 1, len(s.Runs))
		assert.Equal(t, 0, len(s.byRunner))
		assert.Equal(t, 0, len(s.byHost))
	})
	t.Run("single run", func(t *testing.T) {
		runs := []*Run{
			{
				Runners: []Talent{{Name: "amazing"}},
				Hosts:   []Talent{{Name: "wonderful"}},
			},
		}
		s := NewScheduleFrom(runs)
		assert.Equal(t, 1, len(s.Runs))
		assert.Equal(t, 1, len(s.byRunner))
		assert.Equal(t, 1, len(s.byHost))
		assert.Equal(t, 1, len(s.byRunner["amazing"]))
		assert.Equal(t, 1, len(s.byHost["wonderful"]))
	})
	t.Run("multiple runs", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assert.Equal(t, 5, len(s.Runs))
		assert.Equal(t, 2, len(s.byRunner))
		assert.Equal(t, 2, len(s.byHost))
		assert.Equal(t, 3, len(s.byRunner["amazing"]))
		assert.Equal(t, 3, len(s.byHost["wonderful"]))
		assert.Equal(t, 1, len(s.byRunner["fantastic"]))
		assert.Equal(t, 1, len(s.byHost["awesome"]))
	})
}

func TestForEntity(t *testing.T) {
	s := NewScheduleFrom(testRuns)
	t.Run("runner: empty", func(t *testing.T) {
		assert.Equal(t, nil, s.forEntity("runner", " "))
	})
	t.Run("host: empty", func(t *testing.T) {
		assert.Equal(t, nil, s.forEntity("host", " "))
	})
	t.Run("unknown", func(t *testing.T) {
		assert.Equal(t, nil, s.forEntity("runner", "zz"))
	})
	t.Run("exact match", func(t *testing.T) {
		assert.Equal(t, 3, len(s.forEntity("runner", "amazing").Runs))
	})
	t.Run("exact match with extra spacing", func(t *testing.T) {
		assert.Equal(t, 3, len(s.forEntity("runner", " amazing  ").Runs))
	})
	t.Run("patial match single runner", func(t *testing.T) {
		assert.Equal(t, 3, len(s.forEntity("runner", "maz").Runs))
	})
	t.Run("patial match multiple runners", func(t *testing.T) {
		assert.Equal(t, 4, len(s.forEntity("runner", "a").Runs))
	})
	t.Run("unknown kind", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("unknown kind must result in a panic")
			}
		}()
		s.forEntity("derp", "a")
	})

	t.Run("ForRunner", func(t *testing.T) {
		assert.Equal(t, 4, len(s.ForRunner("a").Runs))
	})
	t.Run("ForHost", func(t *testing.T) {
		assert.Equal(t, 1, len(s.ForHost("a").Runs))
	})
}

func TestNextRun(t *testing.T) {
	t.Run("with runs in the future", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		e := s.NextRun(time.Now().UTC())
		assert.Equal(t, "Game 4", e.Title)
	})
	t.Run("with only runs in the past", func(t *testing.T) {
		s := NewScheduleFrom([]*Run{
			{
				Title:   "Game 1",
				Runners: []Talent{{Name: "amazing"}},
				Hosts:   []Talent{{Name: "wonderful"}},
				Start:   time.Now().Add(-10 * time.Minute).UTC(),
			},
		})
		e := s.NextRun(time.Now().UTC())
		assert.Equal(t, nil, e)
	})
}

func TestForTitle(t *testing.T) {
	t.Run("empty title", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assert.Equal(t, nil, s.ForTitle(""))
	})
	t.Run("unknown title", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assert.Equal(t, nil, s.ForTitle(""))
	})
	t.Run("matching title", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assert.Equal(t, 1, len(s.ForTitle("4 ").Runs))
	})
	t.Run("matching multiple titles", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assert.Equal(t, 5, len(s.ForTitle(" ga ").Runs))
	})
}
