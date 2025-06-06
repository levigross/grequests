// Package grequests implements a friendly API over Go's existing net/http library
package grequests

import "context"

// Get takes 2 parameters and returns a Response Struct. These two options are:
//  1. A URL
//  2. A set of options for the request
func Get(ctx context.Context, url string, options ...Option) (*Response, error) {
	return Request(ctx, "GET", url, options...)
}

// Put takes 2 parameters and returns a Response struct. These two options are:
//  1. A URL
//  2. A set of options for the request
//
// If you do not intend to use the `RequestOptions` you can just pass nil
func Put(ctx context.Context, url string, options ...Option) (*Response, error) {
	return Request(ctx, "PUT", url, options...)
}

// Patch takes 2 parameters and returns a Response struct. These two options are:
//  1. A URL
//  2. A set of options for the request
//
// If you do not intend to use the `RequestOptions` you can just pass nil
func Patch(ctx context.Context, url string, options ...Option) (*Response, error) {
	return Request(ctx, "PATCH", url, options...)
}

// Delete takes 2 parameters and returns a Response struct. These two options are:
//  1. A URL
//  2. A set of options for the request
//
// If you do not intend to use the `RequestOptions` you can just pass nil
func Delete(ctx context.Context, url string, options ...Option) (*Response, error) {
	return Request(ctx, "DELETE", url, options...)
}

// Post takes 2 parameters and returns a Response channel. These two options are:
//  1. A URL
//  2. A set of options for the request
//
// If you do not intend to use the `RequestOptions` you can just pass nil
func Post(ctx context.Context, url string, options ...Option) (*Response, error) {
	return Request(ctx, "POST", url, options...)
}

// Head takes 2 parameters and returns a Response channel. These two options are:
//  1. A URL
//  2. A set of options for the request
//
// If you do not intend to use the `RequestOptions` you can just pass nil
func Head(ctx context.Context, url string, options ...Option) (*Response, error) {
	return Request(ctx, "HEAD", url, options...)
}

// Options takes 2 parameters and returns a Response struct. These two options are:
//  1. A URL
//  2. A set of options for the request
//
// If you do not intend to use the `RequestOptions` you can just pass nil
func Options(ctx context.Context, url string, options ...Option) (*Response, error) {
	return Request(ctx, "OPTIONS", url, options...)
}

// Request takes 3 parameters and returns a Response Struct. These three options are:
//  1. A verb
//  2. A URL
//  3. A set of options for the request
//
// If you do not intend to use the `RequestOptions` you can just pass nil
func Request(ctx context.Context, verb, url string, options ...Option) (*Response, error) {
	ro := &RequestOptions{}
	for _, opt := range options {
		opt.Apply(ro)
	}
	if ctx != nil {
		ro.Context = ctx
	}
	return DoRegularRequest(verb, url, ro)
}
