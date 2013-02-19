package v0

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

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
	if isOurs(r) {
		resp, err = processRequest(r)
	} else {
		resp, err = http.DefaultTransport.RoundTrip(r)
		if err != nil {
			err = errors.New(err)
			return
		}
	}

	if resp != nil {
		log.Printf("%s %d %s\n", r.Method, resp.StatusCode, r.URL)
	}
	return
}

func isOurs(r *http.Request) bool {
	u := r.URL.Path
	prefixes := []string{
		"/components/",
		"/scripts/",
		"/styles/",
		"/favicon.ico",
	}
	for _, prefix := range prefixes {
		if len(u) >= len(prefix) && u[:len(prefix)] == prefix {
			return true
		}
	}

	return false
}

func processRequest(r *http.Request) (*http.Response, error) {
	body, err := readFile(filepath.Join("client", "app", r.URL.Path))
	if err != nil {
		return nil, err
	}
	respBody := bytes.NewReader(body)

	return &http.Response{
		StatusCode:    http.StatusOK,
		Status:        http.StatusText(http.StatusOK),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		Body:          ioutil.NopCloser(respBody),
		ContentLength: int64(respBody.Len()),
		Request:       r,
	}, nil
}

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, errors.New(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, f); err != nil {
		return nil, errors.New(err)
	}

	return buf.Bytes(), nil
}
