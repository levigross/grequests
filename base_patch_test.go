package grequests

import (
	"testing"
)

func TestBasicPatchRequest(t *testing.T) {
	resp := Patch("http://httpbin.org/patch", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}
}

func TestBasicAsyncPatchRequest(t *testing.T) {
	resp := <-PatchAsync("http://httpbin.org/patch", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}
}
