// Package grequests implements a friendly API over Go's existing net/http library
package grequests // import "github.com/levigross/grequests"

// Get takes 2 parameters and returns a Response Struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Get(url string, ro *RequestOptions) (*Response, error) {
	return doRequest("GET", url, ro)
}

// GetAsync takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func GetAsync(url string, ro *RequestOptions) chan *Response {
	return doAsyncRequest("GET", url, ro)
}

// Put takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Put(url string, ro *RequestOptions) (*Response, error) {
	return doRequest("PUT", url, ro)
}

// PutAsync takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func PutAsync(url string, ro *RequestOptions) chan *Response {
	return doAsyncRequest("PUT", url, ro)
}

// Patch takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Patch(url string, ro *RequestOptions) (*Response, error) {
	return doRequest("PATCH", url, ro)
}

// PatchAsync takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func PatchAsync(url string, ro *RequestOptions) chan *Response {
	return doAsyncRequest("PATCH", url, ro)
}

// Delete takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Delete(url string, ro *RequestOptions) (*Response, error) {
	return doRequest("DELETE", url, ro)
}

// DeleteAsync takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func DeleteAsync(url string, ro *RequestOptions) chan *Response {
	return doAsyncRequest("DELETE", url, ro)
}

// Post takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Post(url string, ro *RequestOptions) (*Response, error) {
	return doRequest("POST", url, ro)
}

// PostAsync takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func PostAsync(url string, ro *RequestOptions) chan *Response {
	return doAsyncRequest("POST", url, ro)
}

// Head takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Head(url string, ro *RequestOptions) (*Response, error) {
	return doRequest("HEAD", url, ro)
}

// HeadAsync takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func HeadAsync(url string, ro *RequestOptions) chan *Response {
	return doAsyncRequest("HEAD", url, ro)
}

// Options takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Options(url string, ro *RequestOptions) (*Response, error) {
	return doRequest("OPTIONS", url, ro)
}

// OptionsAsync takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func OptionsAsync(url string, ro *RequestOptions) chan *Response {
	return doAsyncRequest("OPTIONS", url, ro)
}
