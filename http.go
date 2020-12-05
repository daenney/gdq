package gdq

import (
	"net"
	"net/http"
	"time"
)

var defaultTrasnport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          5,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	MaxIdleConnsPerHost:   2,
}

type transport struct{}

func (*transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "gdqbot (+https://github.com/daenney/gdq")
	return defaultTrasnport.RoundTrip(req)
}

func newHTTPClient() *http.Client {
	return &http.Client{Transport: &transport{}}
}
