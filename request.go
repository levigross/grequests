package grequests

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RequestOptions is the location that of where the data
type RequestOptions struct {

	// Data is a map of key values that will eventually convert into the query string of a GET request or the
	// body of a POST request. Items can be passed in as an interface (which makes the map easier to construct)
	Data map[string]string

	// Params is a map of query strings that may be used within a GET request
	Params map[string]string

	// Files is where you can include files to upload. The use of this data structure is limited to POST requests
	Files map[string]io.ReadCloser

	// If you want to add custom HTTP headers to the request, this is your friend
	Headers map[string]string

	// InsecureSkipVerify is a flag that specifies if we should validate the server's TLS certificate. It should be noted that
	// Go's TLS verify mechanism doesn't validate if a certificate has been revoked
	InsecureSkipVerify bool

	// UserAgent allows you to set an arbitrary custom user agent
	UserAgent string

	// Auth allows you to specify a user name and password that you wish to use when requesting
	// the URL. It will use basic HTTP authentication formatting the username and password in base64
	// the format is []string{username, password}
	Auth []string
}

// buildRequest is where most of the magic happens for request processing
func buildRequest(httpMethod, url string, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = &RequestOptions{}
	}
	// Create our own HTTP client
	httpClient := buildHTTPClient(ro)
	// Build our URL
	if ro.Params != nil {
		url = buildUrlParams(url, ro.Params)
	}

	// Build the request
	var (
		req *http.Request
		err error
	)

	if IsIdempotentMethod(httpMethod) {
		req, err = http.NewRequest(httpMethod, url, nil)
	} else {
		req, err = buildNonIdempotentRequest(httpMethod, url, ro)
	}

	if err != nil {
		return nil, err
	}

	// Do we need to add any HTTP headers or Basic Auth?
	addHTTPHeaders(ro, req)

	return httpClient.Do(req)
}

func buildNonIdempotentRequest(httpMethod, userUrl string, ro *RequestOptions) (*http.Request, error) {
	if httpMethod == "POST" {
		return buildPostRequest(httpMethod, userUrl, ro)
	}

	return nil, nil // Placeholder

}

func buildPostRequest(httpMethod, userUrl string, ro *RequestOptions) (*http.Request, error) {
	if len(ro.Files) == 0 {
		return createBasicPostRequest(httpMethod, userUrl, ro)
	}

	return nil, nil // Placeholder
}

func createBasicPostRequest(httpMethod, userUrl string, ro *RequestOptions) (*http.Request, error) {
	req, err := http.NewRequest(httpMethod, userUrl, strings.NewReader(encodePostValues(ro.Data)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}
func encodePostValues(postValues map[string]string) string {
	urlValues := &url.Values{}
	for key, value := range postValues {
		urlValues.Set(key, value)
	}
	return urlValues.Encode()
}

func buildHTTPClient(ro *RequestOptions) *http.Client {
	httpClient := &http.Client{}

	if ro.InsecureSkipVerify == true {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return httpClient

}

// buildUrlParams returns a URL with all of the params
// Note: This function will override current URL params if they contradict what is provided in the map
func buildUrlParams(userUrl string, params map[string]string) string {
	parsedUrl, err := url.Parse(userUrl)

	if err != nil {
		return userUrl
	}

	parsedQuery, err := url.ParseQuery(parsedUrl.RawQuery)

	for key, value := range params {
		parsedQuery.Set(key, value)
	}

	return strings.Join(
		[]string{strings.Replace(parsedUrl.String(),
			"?"+parsedUrl.RawQuery, "", -1),
			parsedQuery.Encode()},
		"?")
}

// addHTTPHeaders adds any additional HTTP headers that need to be added are added here including:
// 1. Custom User agent
// 2. Authorization Headers
// 3. Any other header requested
func addHTTPHeaders(ro *RequestOptions, req *http.Request) {
	for key, value := range ro.Headers {
		req.Header.Set(key, value)
	}

	if ro.UserAgent != "" {
		req.Header.Set("User-Agent", ro.UserAgent)
	}

	if ro.Auth != nil {
		req.SetBasicAuth(ro.Auth[0], ro.Auth[1])
	}
}
