package gdq

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testRuns = []*Run{
	{
		Title:   "Game 1",
		Runners: Runners{{Handle: "amazing"}},
		Hosts:   []string{"wonderful"},
		Start:   time.Now().Add(-10 * time.Minute).UTC(),
	},
	{
		Title:   "Game 2",
		Runners: Runners{{Handle: "amazing"}},
	},
	{
		Title: "Game 3",
		Hosts: []string{"wonderful"},
	},
	{
		Title:   "Game 4",
		Runners: Runners{{Handle: "amazing"}},
		Hosts:   []string{"awesome"},
		Start:   time.Now().Add(10 * time.Minute).UTC(),
	},
	{
		Title:   "Game 5",
		Runners: Runners{{Handle: "fantastic"}},
		Hosts:   []string{"wonderful"},
		Start:   time.Now().Add(30 * time.Minute).UTC(),
	},
}

func TestNewScheduleFrom(t *testing.T) {
	t.Run("no runs", func(t *testing.T) {
		runs := []*Run{}
		s := NewScheduleFrom(runs)
		assert.Equal(t, 0, len(s.Runs))
		assert.Equal(t, 0, len(s.byRunner))
		assert.Equal(t, 0, len(s.byHost))
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
				Runners: Runners{{Handle: "amazing"}},
				Hosts:   []string{"wonderful"},
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
		assert.Equal(t, 0, len(s.forEntity("runner", " ").Runs))
	})
	t.Run("host: empty", func(t *testing.T) {
		assert.Equal(t, 0, len(s.forEntity("host", " ").Runs))
	})
	t.Run("unknown", func(t *testing.T) {
		assert.Equal(t, 0, len(s.forEntity("runner", "zz").Runs))
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
		e := s.NextRun()
		assert.Equal(t, "Game 4", e.Title)
	})
	t.Run("with only runs in the past", func(t *testing.T) {
		s := NewScheduleFrom([]*Run{
			{
				Title:   "Game 1",
				Runners: Runners{{Handle: "amazing"}},
				Hosts:   []string{"wonderful"},
				Start:   time.Now().Add(-10 * time.Minute).UTC(),
			},
		})
		e := s.NextRun()
		assert.Nil(t, e)
	})
}

func TestForTitle(t *testing.T) {
	t.Run("empty title", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assert.Equal(t, 0, len(s.ForTitle("").Runs))
	})
	t.Run("unknown title", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assert.Equal(t, 0, len(s.ForTitle("x").Runs))
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

func newTestMux(t *testing.T) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/hosts/34", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"model": "tracker.hostslot", "pk": 1, "fields": {"start_run": 1000, "end_run": 1000, "name": "host 1"}}, {"model": "tracker.hostslot", "pk": 2, "fields": {"start_run": 1001, "end_run": 1001, "name": "host 2"}}, {"model": "tracker.hostslot", "pk": 3, "fields": {"start_run": 1003, "end_run": 1003, "name": "host 2"}}]`)
	})
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query()["type"][0] {
		case "run":
			fmt.Fprintln(w, `[{"model":"tracker.speedrun","pk":1000,"fields":{"event":1000,"name":"Game 1","display_name":"Game 1","twitch_name":"","deprecated_runners":"Runner 1","console":"","commentators":"","description":"","starttime":"2020-12-10T12:30:00Z","endtime":"2020-12-10T13:00:00Z","order":1,"run_time":"0:20:00","setup_time":"0:10:00","coop":false,"category":"Any%","release_year":null,"giantbomb_id":null,"runners":[1],"canonical_url":"https://gamesdonequick.com/tracker/run/1000","public":"Pre-Show Intro (event_id: 1000)"}},{"model":"tracker.speedrun","pk":1001,"fields":{"event":1000,"name":"Game 2","display_name":"Game 2","twitch_name":"","deprecated_runners":"Runner 2","console":"PC","commentators":"","description":"","starttime":"2020-12-10T13:00:00Z","endtime":"2020-12-10T14:00:00Z","order":2,"run_time":"0:45:00","setup_time":"0:15:00","coop":false,"category":"Glitchless","release_year":2009,"giantbomb_id":null,"runners":[2],"canonical_url":"https://gamesdonequick.com/tracker/run/1001","public":"Game 2 (event_id: 1000)"}},{"model":"tracker.speedrun","pk":1002,"fields":{"event":1000,"name":"Game 3","display_name":"Game 3","twitch_name":"","deprecated_runners":"Runner 3, Runner 4","console":"SNES","commentators":"","description":"","starttime":"2020-12-10T14:00:00Z","endtime":"2020-12-10T15:10:00Z","order":3,"run_time":"1:00:00","setup_time":"0:10:00","coop":true,"category":"Any%","release_year":1994,"giantbomb_id":null,"runners":[3,4],"canonical_url":"https://gamesdonequick.com/tracker/run/1002","public":"Game 4 (event_id: 1000)"}}]`)
		case "runner":
			fmt.Fprintln(w, `[{"model":"tracker.runner","pk":1,"fields":{"name":"runner1","stream":"https://www.twitch.tv/runner1","twitter":"","youtube":"","platform":"TWITCH","pronouns":"","donor":null,"public":"runner1"}},{"model":"tracker.runner","pk":2,"fields":{"name":"runner2","stream":"http://www.twitch.tv/runner2","twitter":"runner2","youtube":"","platform":"TWITCH","pronouns":"","donor":null,"public":"runner2"}},{"model":"tracker.runner","pk":3,"fields":{"name":"runner3","stream":"https://www.twitch.tv/runner2","twitter":"runner3","youtube":"","platform":"TWITCH","pronouns":"","donor":null,"public":"runner3"}},{"model":"tracker.runner","pk":4,"fields":{"name":"runner4","stream":"https://www.twitch.tv/runner4","twitter":"runner4","youtube":"https://www.youtube.com/runner4","platform":"TWITCH","pronouns":"","donor":null,"public":"runner4"}}]`)
		default:
			t.Fatalf("%s %s", r.URL.Path, r.URL.RawQuery)
		}
	})
	return mux
}

func TestGetSchedule(t *testing.T) {
	t.Run("with schedule", func(t *testing.T) {
		ts := httptest.NewServer(newTestMux(t))
		defer ts.Close()

		c := New(context.Background(), http.DefaultClient)
		c.base = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

		s, err := c.Schedule(&AGDQ2021)
		assert.Equal(t, err, nil)

		assert.Equal(t, len(s.Runs), 3)
		assert.Equal(t, len(s.byHost), 2)
		assert.Equal(t, len(s.byRunner), 4)
	})
	t.Run("with bad request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{"error": "Malformed Parameters", "exception": "'Missing parameter: type'"}`)
		}))
		defer ts.Close()
		c := New(context.Background(), http.DefaultClient)
		c.base = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

		_, err := c.Schedule(&AGDQ2021)
		assert.Error(t, err)
	})
}
