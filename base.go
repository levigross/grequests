// Package grequests implements a friendly API over Go's existing net/http library
package grequests // import "github.com/levigross/grequests"

// Get takes 2 parameters and returns a Response channel. These two options are:
// 1. A URL
// 2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Get(url string, ro *RequestOptions) chan *Response {
	return doRequest("GET", url, ro)
}

func Put(url string, ro *RequestOptions) chan *Response { return nil }

func Post(url string, ro *RequestOptions) chan *Response {
	return doRequest("POST", url, ro)
}
func Delete(url string, ro *RequestOptions) chan *Response { return nil }

func Head(url string, ro *RequestOptions) chan *Response {
	return doRequest("HEAD", url, ro)
}

func Options(url string, ro *RequestOptions) chan *Response {
	return doRequest("OPTIONS", url, ro)
}
