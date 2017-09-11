package gusha

import (
	"fmt"
)

func ExampleURL() {
	url := URL{
		Method:    "GET",
		Path:      "/",
		Params:    map[string]interface{}{"token": "{{ .APIKey }}"},
		CacheBust: true,
	}

	fmt.Println(url.String())
	// Output: GET /?token={{ .APIKey }}&cache_bust=true
}
