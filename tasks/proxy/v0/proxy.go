package v0

import (
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

func (p *Proxy) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, errors.New(err)
	}

	log.Printf("%s %d %s\n", r.Method, resp.StatusCode, r.URL)

	return resp, nil
}

/*
// Creates a new response based on an error.
func errResponse(r *http.Request, err error) *http.Response {
	// Log the error to the console
	log.Println("ERROR: %s", err)

	// Create the response using the error as the body
	resp := fmt.Sprintf("%s", err)
	return &http.Response{
		Status:        "500 Internal Server Error",
		StatusCode:    500,
		Proto:         r.Proto,
		ProtoMajor:    r.ProtoMajor,
		ProtoMinor:    r.ProtoMinor,
		Body:          ioutil.NopCloser(bytes.NewBufferString(resp)),
		ContentLength: int64(len(resp)),
	}
}

*/
