package grequests

import (
	"testing"
)

func TestBasicOPTIONSRequest(t *testing.T) {
	resp := Options("http://httpbin.org/get", nil)
	if resp.Error != nil {
		t.Error("Unable to make OPTIONS request: ", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Options request did not return success: ", resp.StatusCode)
	}

	if resp.Header.Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, PATCH, OPTIONS" {
		t.Error("Access-Control-Allow-Methods Type Header is unexpected: ", resp.Header)
	}
}

func TestBasicAsyncOPTIONSRequest(t *testing.T) {
	resp := <-OptionsAsync("http://httpbin.org/get", nil)
	if resp.Error != nil {
		t.Error("Unable to make OPTIONS request: ", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Options request did not return success: ", resp.StatusCode)
	}

	if resp.Header.Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, PATCH, OPTIONS" {
		t.Error("Access-Control-Allow-Methods Type Header is unexpected: ", resp.Header)
	}
}
