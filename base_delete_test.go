package grequests

import (
	"testing"
)

func TestBasicDeleteRequest(t *testing.T) {
	resp := Delete("http://httpbin.org/get", nil)
	if resp != nil {
		t.Error("It is time to start writing the DELETE tests")
	}
}
