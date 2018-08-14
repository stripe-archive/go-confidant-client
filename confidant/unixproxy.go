package confidant

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

type Transport struct {
	shadow http.Transport
}

func NewUnixProxyTransport(path string) *Transport {
	const ResponseHeaderTimeout = 30 * time.Second
	const ExpectContinueTimeout = 10 * time.Second
	const DisableKeepAlives = true

	dial := func(network, addr string) (net.Conn, error) {
		return net.Dial("unix", path)
	}

	shadow := http.Transport{
		Dial:                  dial,
		DialTLS:               dial,
		DisableKeepAlives:     DisableKeepAlives,
		ResponseHeaderTimeout: ResponseHeaderTimeout,
		ExpectContinueTimeout: ExpectContinueTimeout,
	}

	return &Transport{shadow}
}

// UnixProxy returns an HTTP client that proxies through a unix socket with the provided path.
func UnixProxy(path string) *Transport {
	return NewUnixProxyTransport(os.ExpandEnv(path))
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	roundTripRequest := *req
	roundTripURL := *req.URL
	roundTripRequest.URL = &roundTripURL
	roundTripRequest.URL.Opaque = fmt.Sprintf("//%s%s", req.URL.Host, req.URL.EscapedPath())
	return t.shadow.RoundTrip(&roundTripRequest)
}
