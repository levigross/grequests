package grequests

import (
	"testing"
)

func TestBasicHeadRequest(t *testing.T) {
	resp := Head("http://httpbin.org/get", nil)
	if resp.Error != nil {
		t.Error("Unable to make HEAD request: ", resp.Error)
	}

	if resp.Ok != true {
		t.Error("HEAD request did not return success: ", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Error("Content Type Header is unexpected: ", resp.Header.Get("Content-Type"))
	}
}

func TestBasicHeadRequestNoContent(t *testing.T) {
	resp := Head("http://httpbin.org/bytes/0", nil)
	if resp.Error != nil {
		t.Error("Unable to make HEAD request: ", resp.Error)
	}

	if resp.Ok != true {
		t.Error("HEAD request did not return success: ", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/octet-stream" {
		t.Error("Content Type Header is unexpected: ", resp.Header.Get("Content-Type"))
	}

	if resp.Bytes() != nil {
		t.Error("Somehow byte buffer is working now (bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("Somehow byte buffer is working now (bytes)", resp.String())
	}
}

func TestBasicHeadAsyncRequest(t *testing.T) {
	resp := <-HeadAsync("http://httpbin.org/get", nil)
	if resp.Error != nil {
		t.Error("Unable to make HEAD request: ", resp.Error)
	}

	if resp.Ok != true {
		t.Error("HEAD request did not return success: ", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Error("Content Type Header is unexpected: ", resp.Header.Get("Content-Type"))
	}
}
