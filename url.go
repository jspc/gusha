package gusha

import (
	"fmt"
	"strings"
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
	params := []string{}
	for k, v := range u.Params {
		params = append(params, fmt.Sprintf("%s=%v", k, v))
	}

	if u.CacheBust {
		params = append(params, "cache_bust=true")
	}

	return fmt.Sprintf("%s %s?%s", u.Method, u.Path, strings.Join(params, "&"))
}
