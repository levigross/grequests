// Package grequests implements a friendly API over Go's existing net/http library
package grequests // import "github.com/levigross/grequests"

// GET takes 2 parameters and returns a Response channel. These two options are:
// 1. A URL
// 2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Get(url string, ro *RequestOptions) chan *Response {
	responseChan := make(chan *Response)
	go func() {
		responseChan <- buildResponse(buildRequest("GET", url, ro))
	}()
	return responseChan
}

func Put(url string, ro *RequestOptions) chan *Response { return nil }
func Post(url string, ro *RequestOptions) chan *Response {
	responseChan := make(chan *Response)
	go func() {
		responseChan <- buildResponse(buildRequest("POST", url, ro))
	}()
	return responseChan
}
func Delete(url string, ro *RequestOptions) chan *Response { return nil }
func Head(url string, ro *RequestOptions) chan *Response {
	responseChan := make(chan *Response)
	go func() {
		responseChan <- buildResponse(buildRequest("HEAD", url, ro))
	}()
	return responseChan
}
func Options(url string, ro *RequestOptions) chan *Response {
	responseChan := make(chan *Response)
	go func() {
		responseChan <- buildResponse(buildRequest("OPTIONS", url, ro))
	}()
	return responseChan
}
