package main

import (
	"fmt"
	"net/http"
	"testing"
)

var (
	emptyParams  = map[string]interface{}{}
	fooBarParams = map[string]interface{}{"foo": "bar"}
)

type testHTTPClient struct {
	url string
}

func (thc *testHTTPClient) Do(r *http.Request) (*http.Response, error) {
	thc.url = r.URL.String()
	return &http.Response{}, nil
}

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

func TestURL(t *testing.T) {
	for _, domain := range []string{"", "https://example.com", "https://maxfactor.co.uk"} {

		for _, test := range []struct {
			name         string
			path         string
			params       map[string]interface{}
			cacheBust    bool
			expectString string
		}{
			{fmt.Sprintf("Request on domain %q, with no params, no cachebust", domain), "/", emptyParams, false, "/?"},
			{fmt.Sprintf("Request on domain %q, with no params, with cachebust", domain), "/", emptyParams, true, "/?cache_bust=true"},
			{fmt.Sprintf("Request on domain %q, with params, no cachebust", domain), "/", fooBarParams, false, "/?foo=bar"},
			{fmt.Sprintf("Request on domain %q, with params, with cachebust", domain), "/", fooBarParams, true, "/?foo=bar&cache_bust=true"},
		} {
			t.Run(test.name, func(t *testing.T) {
				u := URL{
					Method:    "GET",
					Path:      test.path,
					Params:    test.params,
					CacheBust: test.cacheBust,
				}

				uString := u.urlString(domain, "true")
				fqExpectString := fmt.Sprintf("%s%s", domain, test.expectString)

				if fqExpectString != uString {
					fmt.Errorf("expected %q, received %q", fqExpectString, uString)
				}
			})
		}

	}
}

func TestDo(t *testing.T) {
	for _, test := range []struct {
		name      string
		domain    string
		path      string
		params    map[string]interface{}
		cacheBust bool
	}{
		{"Request with no params, no cachebust", "coty.com", "/about-us", emptyParams, false},
		{"Request with no params, with cachebust", "coty.com", "/about-us", emptyParams, true},
		{"Request with params, no cachebust", "coty.com", "/about-us", fooBarParams, false},
		{"Request with params, with cachebust", "coty.com", "/about-us", fooBarParams, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			u := URL{
				Method:    "GET",
				Path:      test.path,
				Params:    test.params,
				CacheBust: test.cacheBust,
			}

			BaseURL = test.domain
			Client = &testHTTPClient{}

			_, err := u.Do()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
