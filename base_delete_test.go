package grequests

import (
	"testing"
)

func TestBasicDeleteRequest(t *testing.T) {
	resp := Delete("http://httpbin.org/delete", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}
}

func TestBasicAsyncDeleteRequest(t *testing.T) {
	resp := <-DeleteAsync("http://httpbin.org/delete", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}
}
