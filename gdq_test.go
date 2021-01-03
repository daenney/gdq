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
	t.Run("with empty body", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[]`)
		}))
		defer ts.Close()

		c := New(context.Background(), http.DefaultClient)
		c.base = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

		_, err := c.Latest()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no known events")
	})
}
