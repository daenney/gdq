package gdq

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

var testEvents = []*Event{
	{
		Title:   "Game 1",
		Runners: []string{"amazing"},
		Hosts:   []string{"wonderful"},
		Start:   time.Now().Add(-10 * time.Minute).UTC(),
	},
	{
		Title:   "Game 2",
		Runners: []string{"amazing"},
	},
	{
		Title: "Game 3",
		Hosts: []string{"wonderful"},
	},
	{
		Title:   "Game 4",
		Runners: []string{"amazing"},
		Hosts:   []string{"awesome"},
		Start:   time.Now().Add(10 * time.Minute).UTC(),
	},
	{
		Title:   "Game 5",
		Runners: []string{"fantastic"},
		Hosts:   []string{"wonderful"},
		Start:   time.Now().Add(30 * time.Minute).UTC(),
	},
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	t.Helper()
	if a == b {
		return
	}
	t.Errorf("Received '%v' (type %v), expected '%v' (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func assertNotNil(t *testing.T, a interface{}) {
	t.Helper()
	if a != nil {
		return
	}
	t.Errorf("Received '%v' (type %v), nil", a, reflect.TypeOf(a))
}

func TestNewScheduleFrom(t *testing.T) {
	t.Run("no events", func(t *testing.T) {
		events := []*Event{}
		s := NewScheduleFrom(events)
		assertEqual(t, len(s.Events), 0)
		assertEqual(t, len(s.byRunner), 0)
		assertEqual(t, len(s.byHost), 0)
	})
	t.Run("single empty event", func(t *testing.T) {
		events := []*Event{{}}
		s := NewScheduleFrom(events)
		assertEqual(t, len(s.Events), 1)
		assertEqual(t, len(s.byRunner), 0)
		assertEqual(t, len(s.byHost), 0)
	})
	t.Run("single event", func(t *testing.T) {
		events := []*Event{
			{
				Runners: []string{"amazing"},
				Hosts:   []string{"wonderful"},
			},
		}
		s := NewScheduleFrom(events)
		assertEqual(t, len(s.Events), 1)
		assertEqual(t, len(s.byRunner), 1)
		assertEqual(t, len(s.byHost), 1)
		assertEqual(t, len(s.byRunner["amazing"]), 1)
		assertEqual(t, len(s.byHost["wonderful"]), 1)
	})
	t.Run("multiple events", func(t *testing.T) {
		s := NewScheduleFrom(testEvents)
		assertEqual(t, len(s.Events), 5)
		assertEqual(t, len(s.byRunner), 2)
		assertEqual(t, len(s.byHost), 2)
		assertEqual(t, len(s.byRunner["amazing"]), 3)
		assertEqual(t, len(s.byHost["wonderful"]), 3)
		assertEqual(t, len(s.byRunner["fantastic"]), 1)
		assertEqual(t, len(s.byHost["awesome"]), 1)
	})
}

func TestForEntity(t *testing.T) {
	s := NewScheduleFrom(testEvents)
	t.Run("runner: empty", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", " ").Events), 0)
	})
	t.Run("host: empty", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("host", " ").Events), 0)
	})
	t.Run("unknown", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", "zz").Events), 0)
	})
	t.Run("exact match", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", "amazing").Events), 3)
	})
	t.Run("exact match with extra spacing", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", " amazing  ").Events), 3)
	})
	t.Run("patial match single runner", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", "maz").Events), 3)
	})
	t.Run("patial match multiple runners", func(t *testing.T) {
		assertEqual(t, len(s.forEntity("runner", "a").Events), 4)
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
		assertEqual(t, len(s.ForRunner("a").Events), 4)
	})
	t.Run("ForHost", func(t *testing.T) {
		assertEqual(t, len(s.ForHost("a").Events), 1)
	})
}

func TestNextEvent(t *testing.T) {
	t.Run("with events in the future", func(t *testing.T) {
		s := NewScheduleFrom(testEvents)
		e := s.NextEvent()
		assertEqual(t, e.Title, "Game 4")
	})
	t.Run("with only events in the past", func(t *testing.T) {
		s := NewScheduleFrom([]*Event{
			{
				Title:   "Game 1",
				Runners: []string{"amazing"},
				Hosts:   []string{"wonderful"},
				Start:   time.Now().Add(-10 * time.Minute).UTC(),
			},
		})
		e := s.NextEvent()
		var empty *Event
		assertEqual(t, e, empty)
	})
}

func TestForTitle(t *testing.T) {
	t.Run("empty title", func(t *testing.T) {
		s := NewScheduleFrom(testEvents)
		assertEqual(t, len(s.ForTitle("").Events), 0)
	})
	t.Run("unknown title", func(t *testing.T) {
		s := NewScheduleFrom(testEvents)
		assertEqual(t, len(s.ForTitle("x").Events), 0)
	})
	t.Run("matching title", func(t *testing.T) {
		s := NewScheduleFrom(testEvents)
		assertEqual(t, len(s.ForTitle("4 ").Events), 1)
	})
	t.Run("matching multiple titles", func(t *testing.T) {
		s := NewScheduleFrom(testEvents)
		assertEqual(t, len(s.ForTitle(" ga ").Events), 5)
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

func TestGetSchedule(t *testing.T) {
	t.Run("without client", func(t *testing.T) {
		_, err := GetSchedule(Latest, nil)
		assertNotNil(t, err)
		if !strings.Contains(err.Error(), "missing") {
			t.Errorf("Got error: %s", err)
		}
	})
	t.Run("URL construction", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			assertEqual(t, req.URL.String(), fmt.Sprintf("%s/0", scheduleURL))
			return newRespWithBody(``)
		})
		_, _ = GetSchedule(Latest, client)
	})
	t.Run("empty response", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			return newRespWithBody(``)
		})

		_, err := GetSchedule(Latest, client)
		assertNotNil(t, err)
		if !errors.Is(err, io.EOF) {
			t.Errorf("Got %v, expected %v", err, io.EOF)
		}
	})
	t.Run("missing runtable", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			return newRespWithBody(`<html></html>`)
		})

		_, err := GetSchedule(Latest, client)
		assertNotNil(t, err)
		if !errors.Is(err, ErrMissingSchedule) {
			t.Errorf("Got %v, expected %v", err, ErrMissingSchedule)
		}
	})
	t.Run("missing runtable body", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			return newRespWithBody(`<html><table id="runTable"></table></html>`)
		})

		_, err := GetSchedule(Latest, client)
		assertNotNil(t, err)
		if !errors.Is(err, ErrMissingSchedule) {
			t.Errorf("Got %v, expected %v", err, ErrMissingSchedule)
		}
	})
	t.Run("empty runtable", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			return newRespWithBody(`<html><table id="runTable"><tbody></tbody></table></html>`)
		})

		_, err := GetSchedule(Latest, client)
		assertNotNil(t, err)
		if !errors.Is(err, ErrMissingSchedule) {
			t.Errorf("Got %v, expected %v", err, ErrMissingSchedule)
		}
	})
	t.Run("uneven runs", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			return newRespWithBody(`<html><table id="runTable"><tbody><tr></tr><tr></tr><tr></tr></tbody></table></html>`)
		})

		_, err := GetSchedule(Latest, client)
		assertNotNil(t, err)
		if !errors.Is(err, ErrInvalidSchedule) {
			t.Errorf("Got %v, expected %v", err, ErrInvalidSchedule)
		}
	})
	t.Run("runtable, empty rows", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			return newRespWithBody(`<html><table id="runTable"><tbody><tr></tr><tr></tr></tbody></table></html>`)
		})

		_, err := GetSchedule(Latest, client)
		assertNotNil(t, err)
		if !errors.Is(err, ErrUnexpectedData) {
			t.Errorf("Got %v, expected %v", err, ErrUnexpectedData)
		}
	})
	t.Run("runtable, multiple events", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			return newRespWithBody(`<html><table id="runTable"><tbody><tr>
			<td>2020-12-01T16:00:00Z</td>
			<td>First Event</td>
			<td>First Runner</td>
			<td> <i class="fa fa-clock-o text-gdq-red" aria-hidden="true"></i> 0:10:00 </td>
			</tr>
			<tr>
			<td> <i class="fa fa-clock-o" aria-hidden="true"></i> 0:20:00 </td>
			<td>Any% &mdash; </td>
			<td><i class="fa fa-microphone"></i> First commentator</td>
			</tr>
			<tr>
			<td>2020-12-01T17:00:00Z</td>
			<td>Second&#039;s Game</td>
			<td>Second Runner</td>
			<td> <i class="fa fa-clock-o text-gdq-red" aria-hidden="true"></i> 0:17:00 </td>
			</tr>
			<tr>
			<td> <i class="fa fa-clock-o" aria-hidden="true"></i> 0:46:00 </td>
			<td>Any% Glitch &mdash; PC</td>
			<td><i class="fa fa-microphone"></i> Second commentator</td>
			</tr>
			<tr>
			<td>2020-12-01T18:00:00Z</td>
			<td>Third Game</td>
			<td>Third Runner, Fourth Runner</td>
			<td> <i class="fa fa-clock-o text-gdq-red" aria-hidden="true"></i> 0:25:00 </td>
			</tr>
			<tr>
			<td> <i class="fa fa-clock-o" aria-hidden="true"></i> 0:93:00 </td>
			<td>BBQ &mdash; GBA</td>
			<td><i class="fa fa-microphone"></i> Third commentator, fourthcommentator</td>
			</tr></tbody></table></html>`)
		})

		s, err := GetSchedule(Latest, client)
		assertEqual(t, err, nil)
		assertEqual(t, len(s.Events), 3)
	})
}
