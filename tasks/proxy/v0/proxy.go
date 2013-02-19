package v0

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("proxy", 0, proxy)
}

func proxy(c config.Config, q *registry.Queue) error {
	u, err := url.Parse("http://localhost:8080")
	if err != nil {
		return errors.New(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &Proxy{}

	http.Handle("/", proxy)

	log.Println("Serving app at http://localhost:9810/...")
	if err := http.ListenAndServe(":9810", nil); err != nil {
		return errors.New(err)
	}
	return nil
}

type Proxy struct {
}

func (p *Proxy) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	if p.isOurs(r) {
		resp, err = p.processRequest(r)
	} else {
		resp, err = http.DefaultTransport.RoundTrip(r)
		if err != nil {
			err = errors.New(err)
			return
		}
	}

	log.Printf("%s %d %s\n", r.Method, resp.StatusCode, r.URL)
	return
}

func (p *Proxy) isOurs(r *http.Request) bool {
	u := r.URL.Path
	prefixes := []string{
		"/",
	}
	for _, prefix := range prefixes {
		if u[:len(prefix)] == prefix {
			return true
		}
	}

	return false
}

func (p *Proxy) processRequest(r *http.Request) (*http.Response, error) {
	code := http.StatusOK

	body := []byte("hello")
	respBody := bytes.NewReader(body)

	resp := &http.Response{
		StatusCode:    code,
		Status:        http.StatusText(code),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		Body:          ioutil.NopCloser(respBody),
		ContentLength: int64(respBody.Len()),
		Request:       r,
	}
	return resp, nil
}
