package gdq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func newTestMux(t *testing.T) *http.ServeMux {
	t.Helper()

	var runBuf bytes.Buffer
	runData, err := os.ReadFile("testdata/runs-34.json")
	if err != nil {
		t.Fatalf("got error reading data: %s", err)
	}
	if err := json.Compact(&runBuf, runData); err != nil {
		t.Fatalf("data is not valid JSON: %s", err)
	}

	var ivBuf bytes.Buffer
	ivData, err := os.ReadFile("testdata/interviews-34.json")
	if err != nil {
		t.Fatalf("got error reading data: %s", err)
	}
	if err := json.Compact(&ivBuf, ivData); err != nil {
		t.Fatalf("data is not valid JSON: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/events/34/runs/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(runBuf.Bytes())
	})
	mux.HandleFunc("/events/34/interviews/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(ivBuf.Bytes())
	})
	return mux
}

func TestGetSchedule(t *testing.T) {
	ts := httptest.NewServer(newTestMux(t))
	defer ts.Close()

	c := New(http.DefaultClient)
	c.v2 = fmt.Sprintf("http://%s", ts.Listener.Addr().String())

	s, err := c.Schedule(context.TODO(), AGDQ2021.ID)
	assert.NoError(t, err)
	assert.Equal(t, 157, len(s.Runs))
	assert.Equal(t, 31, len(s.byHost))
	assert.Equal(t, 162, len(s.byRunner))
}

func TestGetWithCtx(t *testing.T) {
	t.Run("with success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `hello`)
		}))
		defer ts.Close()
		ctx := context.Background()
		resp, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
		assert.NoError(t, err)
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
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "failed to unmarshal")
		})
		t.Run("with unexpected JSON body", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `{"test": true}`)
			}))
			defer ts.Close()
			ctx := context.Background()
			_, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "unexpected body")
		})
		t.Run("with JSON error", func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `{"detail": "Malformed some parameter"}`)
			}))
			defer ts.Close()
			ctx := context.Background()
			_, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Malformed some parameter")
		})
	})
	t.Run("with not found", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"detail":"That resource does not exist or you do not have permission to view it."}`)
		}))
		defer ts.Close()
		ctx := context.Background()
		_, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "resource does not exist")
	})
	t.Run("with something else", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, ``)
		}))
		defer ts.Close()
		ctx := context.Background()
		_, err := getWithCtx(ctx, http.DefaultClient, ts.URL)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), http.StatusText(http.StatusBadGateway))
	})
}
