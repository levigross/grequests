// Package grequests implements a friendly API over Go's existing net/http library
package grequests

// Get takes 2 parameters and returns a Response Struct. These two options are:
// 	1. A URL
// 	2. A varible number of options
func Get(url string, options ...Option) (*Response, error) {
	ro := &RequestOptions{}
	for _, opt := range options {
		opt.Apply(ro)
	}
	return DoRegularRequest("GET", url, ro)
}

// Put takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Put(url string, options ...Option) (*Response, error) {
	ro := &RequestOptions{}
	for _, opt := range options {
		opt.Apply(ro)
	}
	return DoRegularRequest("PUT", url, ro)
}

// Patch takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Patch(url string, options ...Option) (*Response, error) {
	ro := &RequestOptions{}
	for _, opt := range options {
		opt.Apply(ro)
	}
	return DoRegularRequest("PATCH", url, ro)
}

// Delete takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Delete(url string, options ...Option) (*Response, error) {
	var ro *RequestOptions
	for _, opt := range options {
		opt.Apply(ro)
	}
	return DoRegularRequest("DELETE", url, ro)
}

// Post takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Post(url string, options ...Option) (*Response, error) {
	ro := &RequestOptions{}
	for _, opt := range options {
		opt.Apply(ro)
	}
	return DoRegularRequest("POST", url, ro)
}

// Head takes 2 parameters and returns a Response channel. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Head(url string, options ...Option) (*Response, error) {
	ro := &RequestOptions{}
	for _, opt := range options {
		opt.Apply(ro)
	}
	return DoRegularRequest("HEAD", url, ro)
}

// Options takes 2 parameters and returns a Response struct. These two options are:
// 	1. A URL
// 	2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Options(url string, options ...Option) (*Response, error) {
	ro := &RequestOptions{}
	for _, opt := range options {
		opt.Apply(ro)
	}
	return DoRegularRequest("OPTIONS", url, ro)
}

// Req takes 3 parameters and returns a Response Struct. These three options are:
//	1. A verb
// 	2. A URL
// 	3. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Req(verb string, url string, options ...Option) (*Response, error) {
	ro := &RequestOptions{}
	for _, opt := range options {
		opt.Apply(ro)
	}
	return DoRegularRequest(verb, url, ro)
}
