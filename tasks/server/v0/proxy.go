package v0

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ernestokarim/cb/config"
)

type proxy struct {
	hosts map[string]string
}

func (p *proxy) RoundTrip(r *http.Request) (*http.Response, error) {
	found := false
	for k, v := range p.hosts {
		if !found && k == r.Host {
			r.Host = v
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("host `%s` not found in mappings: %+v", r.Host, p.hosts)
	}

	// Make the real request
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, fmt.Errorf("roundtrip failed: %s", err)
	}

	// Log the request data
	length := resp.Header.Get("Content-Length")
	var size int64
	if length != "" {
		size, err = strconv.ParseInt(length, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse resp size: %s", err)
		}
	}
	var zero time.Time
	writeLog(r, zero, resp.StatusCode, int(size))

	// Rewrite the location header to the new host if present
	if resp.StatusCode == 302 || resp.StatusCode == 301 {
		location, err := url.Parse(resp.Header.Get("Location"))
		if err != nil {
			return nil, fmt.Errorf("cannot parse the redirect url: %s", err)
		}
		location.Host = fmt.Sprintf("%s:%d", r.Host, *config.Port)
		resp.Header.Set("Location", location.String())
	}

	return resp, nil
}
