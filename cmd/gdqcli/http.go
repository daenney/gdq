package main

import (
	"net/http"
	"time"
)

type transport struct {
	userAgent string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.userAgent != "" {
		req.Header.Set("User-Agent", t.userAgent)
	}
	return http.DefaultTransport.RoundTrip(req)
}

func newHTTPClient(ua string) *http.Client {
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: &transport{userAgent: ua},
	}
}
