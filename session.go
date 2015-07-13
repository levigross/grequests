package grequests

import "net/http"

// Session allows a user to make use of persistent cookies in between
// HTTP requests
type Session struct {

	// HTTPClient is the client that we will use to request the resources
	HTTPClient *http.Client
}

// NewSession returns a session struct which enables can be used to maintain establish a persistent state with the
// server
// This function will set UseCookieJar to true as that is the purpose of using the session
func NewSession(ro *RequestOptions) *Session {
	if ro == nil {
		ro = &RequestOptions{}
	}

	ro.UseCookieJar = true

	return &Session{HTTPClient: BuildHTTPClient(*ro)}
}

// Get takes 2 parameters and returns a Response Struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
// A new session is created by calling NewSession with a request options struct
func (s *Session) Get(url string, ro *RequestOptions) (*Response, error) {
	return doSessionRequest("GET", url, ro, s.HTTPClient)
}

// Put takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
// A new session is created by calling NewSession with a request options struct
func (s *Session) Put(url string, ro *RequestOptions) (*Response, error) {
	return doSessionRequest("PUT", url, ro, s.HTTPClient)
}

// Patch takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
// A new session is created by calling NewSession with a request options struct
func (s *Session) Patch(url string, ro *RequestOptions) (*Response, error) {
	return doSessionRequest("PATCH", url, ro, s.HTTPClient)
}

// Delete takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
// A new session is created by calling NewSession with a request options struct
func (s *Session) Delete(url string, ro *RequestOptions) (*Response, error) {
	return doSessionRequest("DELETE", url, ro, s.HTTPClient)
}

// Post takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
// A new session is created by calling NewSession with a request options struct
func (s *Session) Post(url string, ro *RequestOptions) (*Response, error) {
	return doSessionRequest("POST", url, ro, s.HTTPClient)
}

// Head takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
// A new session is created by calling NewSession with a request options struct
func (s *Session) Head(url string, ro *RequestOptions) (*Response, error) {
	return doSessionRequest("HEAD", url, ro, s.HTTPClient)
}

// Options takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
// A new session is created by calling NewSession with a request options struct
func (s *Session) Options(url string, ro *RequestOptions) (*Response, error) {
	return doSessionRequest("OPTIONS", url, ro, s.HTTPClient)
}

// CloseIdleConnections closes the idle connections that a session client may make use of
func (s *Session) CloseIdleConnections() {
	s.HTTPClient.Transport.(*http.Transport).CloseIdleConnections()
}
