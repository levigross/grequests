package grequests

import (
	"encoding/json"
	"io"
	"net/http"
)

// RequestOptions is the location that of where the data
type RequestOptions struct {
	// Data is a map of key values that will eventually convert into the query string of a GET request or the
	// body of a POST request. Items can be passed in as an interface (which makes the map easier to construct)

	Data map[string]interface{}

	// VerifySSL is a flag that specifies if we should validate the server's TLS certificate. It should be noted that
	// Go's TLS verify mechanism doesn't validate if a certificate has been revoked

	VerifyTLS bool

	// This allows you to set an arbitrary custom user agent

	UserAgent string
}

type Response struct {
	// Ok is a boolean flag that validates that the server returned a 2xx code

	Ok bool

	// This is the Go error flag – if something went wrong within the request, this flag will be set.

	Error error

	// We want to abstract (at least at the moment) the Go http.Response object away. So we are going to make use of it
	// internal but not give the user access

	resp *http.Response
}

func buildResponse(resp *http.Response, err error) *Response {
	// If the connection didn't succeed we just return a blank response
	if err != nil {
		return &Response{Error: err}
	}

	return &Response{
		// If your code is within the 2xx range – the response is considered `Ok`
		Ok:    resp.StatusCode <= 200 && resp.StatusCode <300,
		Error: nil,
		resp:  resp,
	}
}

// JSON is a function that will populate a struct that is provided `userStruct` with the JSON returned within the
// response body
func (r *Response) JSON(userStruct interface{}) error {
	jsonDecoder := json.NewDecoder(r.resp.Body)
	defer r.resp.Body.Close()

	if err := jsonDecoder.Decode(&userStruct); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// GET takes 2 parameters and returns a Response channel. These two options are:
// 1. A URL
// 2. A RequestOptions struct
// If you do not intend to use the `RequestOptions` you can just pass nil
func Get(url string, ro *RequestOptions) chan *Response {
	responseChan := make(chan *Response)
	go func() {
		if ro == nil {
			responseChan <- buildResponse(http.Get(url))
		}
	}()
	return responseChan
}

func Put(url string, ro *RequestOptions) chan *Response     { return nil }
func Post(url string, ro *RequestOptions) chan *Response    { return nil }
func Delete(url string, ro *RequestOptions) chan *Response  { return nil }
func Head(url string, ro *RequestOptions) chan *Response    { return nil }
func Options(url string, ro *RequestOptions) chan *Response { return nil }
