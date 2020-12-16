package main

import (
	"net/http"
)

const userAgent = "gdqcli (+https://github.com/daenney/gdq)"

type transport struct{}

func (*transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", userAgent)
	return http.DefaultTransport.RoundTrip(req)
}

func newHTTPClient() *http.Client {
	return &http.Client{Transport: &transport{}}
}
