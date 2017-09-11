package gusha

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/satori/go.uuid"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

var (
	BaseURL            = "example.com"
	Client  httpClient = &http.Client{}
)

// URLs is a shorthand type for a slice of URL types
type URLs []URL

// URL represents a URL which a gusha agent hits
//
// Note: Cache bust will add the query param `cache_bust`
// and set the value to a new UUID each time
type URL struct {
	Method    string
	Path      string
	Params    map[string]interface{}
	CacheBust bool
}

// String returns a textual representation of a URL
// suitable for logging
func (u URL) String() string {
	return fmt.Sprintf("%s %s", u.Method, u.urlString("", "true"))
}

// Do wraps net/http Client.Do and shadows its return types.
// It exists solely to allow for us to manipulate URLs at runtime
// where Cache Bust behaviour is required
func (u URL) Do() (r *http.Response, err error) {
	pU, err := url.Parse(u.urlString(BaseURL, uuid.NewV4().String()))
	if err != nil {
		return
	}

	return Client.Do(&http.Request{
		Method: u.Method,
		URL:    pU,
	})
}

func (u URL) urlString(base, cacheBuster string) string {
	params := []string{}
	for k, v := range u.Params {
		params = append(params, fmt.Sprintf("%s=%v", k, v))
	}

	if u.CacheBust {
		params = append(params, fmt.Sprintf("cache_bust=%s", cacheBuster))
	}

	return fmt.Sprintf("%s%s?%s", base, u.Path, strings.Join(params, "&"))
}
