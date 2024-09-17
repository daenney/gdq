package gdq

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWithCtx(t *testing.T) {
	t.Run("with success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `hello`)
		}))
		defer ts.Close()
		ctx := context.Background()
		resp, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
		assert.Nil(t, err)
		assert.Equal(t, "hello", string(resp))
	})
	t.Run("with bad request", func(t *testing.T) {
		t.Run("with non-JSON body", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `hello`)
			}))
			defer ts.Close()
			ctx := context.Background()
			_, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
			assert.NotNil(t, err)
			if !strings.Contains(err.Error(), "failed to unmarshal") {
				t.Fatal(err.Error())
			}
		})
		t.Run("with unexpected JSON body", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `{"test": true}`)
			}))
			defer ts.Close()
			ctx := context.Background()
			_, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
			assert.NotNil(t, err)
			if !strings.Contains(err.Error(), "unexpected body") {
				t.Fatal(err.Error())
			}
		})
		t.Run("with JSON error", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `{"error": "Malformed Parameters", "exception": "'Missing parameter: type'"}`)
			}))
			defer ts.Close()
			ctx := context.Background()
			_, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
			assert.NotNil(t, err)
			if !strings.Contains(err.Error(), "Malformed Parameters") {
				t.Fatal(err.Error())
			}
		})
	})
	t.Run("with not found", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, ``)
		}))
		defer ts.Close()
		ctx := context.Background()
		_, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
		assert.NotNil(t, err)
		if !strings.Contains(err.Error(), "resource not found") {
			t.Fatal(err.Error())
		}
	})
	t.Run("with something else", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, ``)
		}))
		defer ts.Close()
		ctx := context.Background()
		_, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
		assert.NotNil(t, err)
		if !strings.Contains(err.Error(), fmt.Sprint(http.StatusBadGateway)) {
			t.Fatal(err.Error())
		}
	})
}

func TestLatest(t *testing.T) {
	t.Run("with invalid body", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `hello`)
		}))
		defer ts.Close()

		c := New(context.Background(), http.DefaultClient)
		c.base = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

		_, err := c.Latest()
		assert.NotNil(t, err)
	})
	t.Run("with events", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"model":"tracker.event","pk":1,"fields":{"short":"event1","name":"Event 1 2020","hashtag":"","use_one_step_screening":true,"receivername":"","targetamount":100,"minimumdonation":1,"paypalemail":"example@example.com","paypalcurrency":"USD","datetime":"2020-05-01T13:00:00Z","timezone":"US/Eastern","locked":true,"allow_donations":true,"canonical_url":"https://gamesdonequick.com/tracker/event/1","public":"Event 1","amount":0,"count":0,"max":0,"avg":0,"allowed_prize_countries":[],"disallowed_prize_regions":[]}},{"model":"tracker.event","pk":2,"fields":{"short":"event2","name":"Event 2 2020","hashtag":"","use_one_step_screening":true,"receivername":"","targetamount":100,"minimumdonation":1,"paypalemail":"example@example.com","paypalcurrency":"USD","datetime":"2020-10-01T13:00:00Z","timezone":"US/Eastern","locked":true,"allow_donations":true,"canonical_url":"https://gamesdonequick.com/tracker/event/2","public":"Event 2 2020","amount":0,"count":0,"max":0,"avg":0,"allowed_prize_countries":[],"disallowed_prize_regions":[]}}]`)
		}))
		defer ts.Close()

		c := New(context.Background(), http.DefaultClient)
		c.base = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

		ev, err := c.Latest()
		assert.Nil(t, err)
		assert.Equal(t, uint(2), ev.ID)
	})
	t.Run("without events", func(t *testing.T) {
		t.Run("with not one event", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, `[]`)
			}))
			defer ts.Close()

			c := New(context.Background(), http.DefaultClient)
			c.base = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

			_, err := c.Latest()
			assert.Error(t, err)
		})
	})
}

func TestDonations(t *testing.T) {
	t.Run("with invalid body", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `hello`)
		}))
		defer ts.Close()

		c := New(context.Background(), http.DefaultClient)
		c.base = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

		_, err := c.Donations(&Event{})
		assert.NotNil(t, err)
	})
	t.Run("with event", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"model":"tracker.event","pk":1,"fields":{"short":"event1","name":"Event 1 2020","hashtag":"","use_one_step_screening":true,"receivername":"","targetamount":100,"minimumdonation":1,"paypalemail":"example@example.com","paypalcurrency":"USD","datetime":"2020-05-01T13:00:00Z","timezone":"US/Eastern","locked":true,"allow_donations":true,"canonical_url":"https://gamesdonequick.com/tracker/event/1","public":"Event 1","amount":101.3,"count":5,"max":10,"avg":3.5,"allowed_prize_countries":[],"disallowed_prize_regions":[]}}]`)
		}))
		defer ts.Close()

		c := New(context.Background(), http.DefaultClient)
		c.base = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

		do, err := c.Donations(&Event{})
		assert.Nil(t, err)
		assert.Equal(t, &Donations{
			Total:   101.3,
			Max:     10,
			Average: 3.5,
			Count:   5,
		}, do)
	})
	t.Run("without events", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `[]`)
		}))
		defer ts.Close()

		c := New(context.Background(), http.DefaultClient)
		c.base = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

		_, err := c.Donations(&Event{})
		assert.Error(t, err)
	})
}
