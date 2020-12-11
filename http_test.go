package gdq

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetUserAgent(t *testing.T) {
	client := newHTTPClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, r.Header["User-Agent"][0], userAgent)
		fmt.Fprintln(w, "")
	}))
	defer ts.Close()

	_, _ = client.Get(ts.URL)
}
