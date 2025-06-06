// Package grequests implements a friendly API over Go's existing net/http library
package grequests

// Get takes 2 parameters and returns a Response Struct. These two options are:
// 	1. A URL
// 	2. A set of options for the request
func Get(url string, options ...Option) (*Response, error) {
	return Request("GET", url, options...)
}

// Put takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A set of options for the request
// If you do not intend to use the `RequestOptions` you can just pass nil
func Put(url string, options ...Option) (*Response, error) {
	return Request("PUT", url, options...)
}

// Patch takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A set of options for the request
// If you do not intend to use the `RequestOptions` you can just pass nil
func Patch(url string, options ...Option) (*Response, error) {
	return Request("PATCH", url, options...)
}

// Delete takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A set of options for the request
// If you do not intend to use the `RequestOptions` you can just pass nil
func Delete(url string, options ...Option) (*Response, error) {
	return Request("DELETE", url, options...)
}

// Post takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A set of options for the request
// If you do not intend to use the `RequestOptions` you can just pass nil
func Post(url string, options ...Option) (*Response, error) {
	return Request("POST", url, options...)
}

// Head takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A set of options for the request
// If you do not intend to use the `RequestOptions` you can just pass nil
func Head(url string, options ...Option) (*Response, error) {
	return Request("HEAD", url, options...)
}

// Options takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A set of options for the request
// If you do not intend to use the `RequestOptions` you can just pass nil
func Options(url string, options ...Option) (*Response, error) {
	return Request("OPTIONS", url, options...)
}

// Request takes 3 parameters and returns a Response Struct. These three options are:
//	1. A verb
// 	2. A URL
// 	3. A set of options for the request
// If you do not intend to use the `RequestOptions` you can just pass nil
func Request(verb string, url string, options ...Option) (*Response, error) {
	ro := &RequestOptions{}
	for _, opt := range options {
		if opt != nil {
			opt.Apply(ro)
		}
	}
	return DoRegularRequest(verb, url, ro)
}
