package gdq

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
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

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	t.Helper()
	t.Log(a == nil)
	t.Log(b == nil)
	if a == nil && b == nil {
		return
	}
	if a == b {
		return
	}
	t.Fatalf("Received '%v' (type %v), expected '%v' (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func assertNotNil(t *testing.T, a interface{}) {
	t.Helper()
	if a != nil {
		return
	}
	t.Fatalf("Received '%v' (type %v), nil", a, reflect.TypeOf(a))
}

func TestNewScheduleFrom(t *testing.T) {
	t.Run("no runs", func(t *testing.T) {
		runs := []*Run{}
		s := NewScheduleFrom(runs)
		assertEqual(t, len(s.Runs), 0)
		assertEqual(t, len(s.byRunner), 0)
		assertEqual(t, len(s.byHost), 0)
	})
	t.Run("single empty run", func(t *testing.T) {
		runs := []*Run{{}}
		s := NewScheduleFrom(runs)
		assertEqual(t, len(s.Runs), 1)
		assertEqual(t, len(s.byRunner), 0)
		assertEqual(t, len(s.byHost), 0)
	})
	t.Run("single run", func(t *testing.T) {
		runs := []*Run{
			{
				Runners: Runners{{Handle: "amazing"}},
				Hosts:   []string{"wonderful"},
			},
		}
		s := NewScheduleFrom(runs)
		assertEqual(t, len(s.Runs), 1)
		assertEqual(t, len(s.byRunner), 1)
		assertEqual(t, len(s.byHost), 1)
		assertEqual(t, len(s.byRunner["amazing"]), 1)
		assertEqual(t, len(s.byHost["wonderful"]), 1)
	})
	t.Run("multiple runs", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assertEqual(t, len(s.Runs), 5)
		assertEqual(t, len(s.byRunner), 2)
		assertEqual(t, len(s.byHost), 2)
		assertEqual(t, len(s.byRunner["amazing"]), 3)
		assertEqual(t, len(s.byHost["wonderful"]), 3)
		assertEqual(t, len(s.byRunner["fantastic"]), 1)
		assertEqual(t, len(s.byHost["awesome"]), 1)
	})
}

func TestForEntity(t *testing.T) {
	s := NewScheduleFrom(testRuns)
	t.Run("runner: empty", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", " ").Runs), 0)
	})
	t.Run("host: empty", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("host", " ").Runs), 0)
	})
	t.Run("unknown", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", "zz").Runs), 0)
	})
	t.Run("exact match", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", "amazing").Runs), 3)
	})
	t.Run("exact match with extra spacing", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", " amazing  ").Runs), 3)
	})
	t.Run("patial match single runner", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", "maz").Runs), 3)
	})
	t.Run("patial match multiple runners", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", "a").Runs), 4)
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
		assertEqual(t, len(s.ForRunner("a").Runs), 4)
	})
	t.Run("ForHost", func(t *testing.T) {
		assertEqual(t, len(s.ForHost("a").Runs), 1)
	})
}

func TestNextRun(t *testing.T) {
	t.Run("with runs in the future", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		e := s.NextRun()
		assertEqual(t, e.Title, "Game 4")
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
		var empty *Run
		assertEqual(t, e, empty)
	})
}

func TestForTitle(t *testing.T) {
	t.Run("empty title", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assertEqual(t, len(s.ForTitle("").Runs), 0)
	})
	t.Run("unknown title", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assertEqual(t, len(s.ForTitle("x").Runs), 0)
	})
	t.Run("matching title", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assertEqual(t, len(s.ForTitle("4 ").Runs), 1)
	})
	t.Run("matching multiple titles", func(t *testing.T) {
		s := NewScheduleFrom(testRuns)
		assertEqual(t, len(s.ForTitle(" ga ").Runs), 5)
	})
}

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

func newRespWithBody(body string) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		Header: http.Header{
			"content-type": []string{"text/html; charset=UTF-8"},
		},
	}
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

		s, err := c.GetSchedule(AGDQ2021)
		assertEqual(t, err, nil)

		assertEqual(t, len(s.Runs), 3)
		assertEqual(t, len(s.byHost), 2)
		assertEqual(t, len(s.byRunner), 4)
	})
	t.Run("with bad request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{"error": "Malformed Parameters", "exception": "'Missing parameter: type'"}`)
		}))
		defer ts.Close()
		c := New(context.Background(), http.DefaultClient)
		c.base = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

		_, err := c.GetSchedule(AGDQ2021)
		assertNotNil(t, err)
	})
	// t.Run("URL construction", func(t *testing.T) {
	// 	client := newTestClient(func(req *http.Request) *http.Response {
	// 		assertEqual(t, req.URL.String(), fmt.Sprintf("%s/0", scheduleURL))
	// 		return newRespWithBody(``)
	// 	})
	// 	_, _ = GetSchedule(AGDQ2021, client)
	// })
	// t.Run("empty response", func(t *testing.T) {
	// 	client := newTestClient(func(req *http.Request) *http.Response {
	// 		return newRespWithBody(``)
	// 	})

	// 	_, err := GetSchedule(AGDQ2021, client)
	// 	assertNotNil(t, err)
	// 	if !errors.Is(err, io.EOF) {
	// 		t.Errorf("Got %v, expected %v", err, io.EOF)
	// 	}
	// })
	// t.Run("missing runtable", func(t *testing.T) {
	// 	client := newTestClient(func(req *http.Request) *http.Response {
	// 		return newRespWithBody(`<html></html>`)
	// 	})

	// 	_, err := GetSchedule(AGDQ2021, client)
	// 	assertNotNil(t, err)
	// 	if !errors.Is(err, ErrMissingSchedule) {
	// 		t.Errorf("Got %v, expected %v", err, ErrMissingSchedule)
	// 	}
	// })
	// t.Run("missing runtable body", func(t *testing.T) {
	// 	client := newTestClient(func(req *http.Request) *http.Response {
	// 		return newRespWithBody(`<html><table id="runTable"></table></html>`)
	// 	})

	// 	_, err := GetSchedule(AGDQ2021, client)
	// 	assertNotNil(t, err)
	// 	if !errors.Is(err, ErrMissingSchedule) {
	// 		t.Errorf("Got %v, expected %v", err, ErrMissingSchedule)
	// 	}
	// })
	// t.Run("empty runtable", func(t *testing.T) {
	// 	client := newTestClient(func(req *http.Request) *http.Response {
	// 		return newRespWithBody(`<html><table id="runTable"><tbody></tbody></table></html>`)
	// 	})

	// 	_, err := GetSchedule(AGDQ2021, client)
	// 	assertNotNil(t, err)
	// 	if !errors.Is(err, ErrMissingSchedule) {
	// 		t.Errorf("Got %v, expected %v", err, ErrMissingSchedule)
	// 	}
	// })
	// t.Run("uneven runs", func(t *testing.T) {
	// 	client := newTestClient(func(req *http.Request) *http.Response {
	// 		return newRespWithBody(`<html><table id="runTable"><tbody><tr></tr><tr></tr><tr></tr></tbody></table></html>`)
	// 	})

	// 	_, err := GetSchedule(AGDQ2021, client)
	// 	assertNotNil(t, err)
	// 	if !errors.Is(err, ErrInvalidSchedule) {
	// 		t.Errorf("Got %v, expected %v", err, ErrInvalidSchedule)
	// 	}
	// })
	// t.Run("runtable, empty rows", func(t *testing.T) {
	// 	client := newTestClient(func(req *http.Request) *http.Response {
	// 		return newRespWithBody(`<html><table id="runTable"><tbody><tr></tr><tr></tr></tbody></table></html>`)
	// 	})

	// 	_, err := GetSchedule(AGDQ2021, client)
	// 	assertNotNil(t, err)
	// 	if !errors.Is(err, ErrUnexpectedData) {
	// 		t.Errorf("Got %v, expected %v", err, ErrUnexpectedData)
	// 	}
	// })
	// t.Run("runtable, multiple runs", func(t *testing.T) {
	// 	client := newTestClient(func(req *http.Request) *http.Response {
	// 		return newRespWithBody(`<html><table id="runTable"><tbody><tr>
	// 		<td>2020-12-01T16:00:00Z</td>
	// 		<td>First Event</td>
	// 		<td>First Runner</td>
	// 		<td> <i class="fa fa-clock-o text-gdq-red" aria-hidden="true"></i> 0:10:00 </td>
	// 		</tr>
	// 		<tr>
	// 		<td> <i class="fa fa-clock-o" aria-hidden="true"></i> 0:20:00 </td>
	// 		<td>Any% &mdash; </td>
	// 		<td><i class="fa fa-microphone"></i> First commentator</td>
	// 		</tr>
	// 		<tr>
	// 		<td>2020-12-01T17:00:00Z</td>
	// 		<td>Second&#039;s Game</td>
	// 		<td>Second Runner</td>
	// 		<td> <i class="fa fa-clock-o text-gdq-red" aria-hidden="true"></i> 0:17:00 </td>
	// 		</tr>
	// 		<tr>
	// 		<td> <i class="fa fa-clock-o" aria-hidden="true"></i> 0:46:00 </td>
	// 		<td>Any% Glitch &mdash; PC</td>
	// 		<td><i class="fa fa-microphone"></i> Second commentator</td>
	// 		</tr>
	// 		<tr>
	// 		<td>2020-12-01T18:00:00Z</td>
	// 		<td>Third Game</td>
	// 		<td>Third Runner, Fourth Runner</td>
	// 		<td> <i class="fa fa-clock-o text-gdq-red" aria-hidden="true"></i> 0:25:00 </td>
	// 		</tr>
	// 		<tr>
	// 		<td> <i class="fa fa-clock-o" aria-hidden="true"></i> 0:93:00 </td>
	// 		<td>BBQ &mdash; GBA</td>
	// 		<td><i class="fa fa-microphone"></i> Third commentator, fourthcommentator</td>
	// 		</tr></tbody></table></html>`)
	// 	})

	// 	s, err := GetSchedule(AGDQ2021, client)
	// 	assertEqual(t, err, nil)
	// 	assertEqual(t, len(s.Runs), 3)
	// })
}
