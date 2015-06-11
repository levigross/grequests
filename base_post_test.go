package grequests

import (
	"testing"
)

type BasicPostResponse struct {
	Args  struct{} `json:"args"`
	Data  string   `json:"data"`
	Files struct{} `json:"files"`
	Form  struct {
		One string `json:"one"`
	} `json:"form"`
	Headers struct {
		Accept         string `json:"Accept"`
		Content_Length string `json:"Content-Length"`
		Content_Type   string `json:"Content-Type"`
		Host           string `json:"Host"`
		User_Agent     string `json:"User-Agent"`
	} `json:"headers"`
	JSON   interface{} `json:"json"`
	Origin string      `json:"origin"`
	URL    string      `json:"url"`
}

func TestBasicPostRequest(t *testing.T) {
	verifyOkPostResponse(<-Post("http://httpbin.org/post",
		&RequestOptions{Data: map[string]string{"One": "Two"}}), t)

}

// verifyResponse will verify the following conditions
// 1. The request didn't return an error
// 2. The response returned an OK (a status code within the 200 range)
// 3. The output can be coerced to JSON (this may change later)
// It should only be run when testing GET request to http://httpbin.org/post expecting JSON
func verifyOkPostResponse(resp *Response, t *testing.T) *BasicPostResponse {
	if resp.Error != nil {
		t.Fatal("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJsonStruct := &BasicPostResponse{}

	err := resp.Json(myJsonStruct)
	if err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJsonStruct.URL != "http://httpbin.org/post" {
		t.Error("For some reason the URL isn't the same", myJsonStruct.URL)
	}

	if myJsonStruct.Headers.Host != "httpbin.org" {
		t.Error("The host header is invalid")
	}

	if myJsonStruct.Form.One != "Two" {
		t.Error("Invalid post response: ", myJsonStruct.Form.One)
	}

	if resp.Bytes() != nil {
		t.Error("JSON decoding did not fully consume the response stream (Bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("JSON decoding did not fully consume the response stream (String)", resp.String())
	}

	if resp.StatusCode != 200 {
		t.Error("Response returned a non-200 code")
	}

	return myJsonStruct
}
