package grequests

import (
	"testing"
)

func TestBasicPutRequest(t *testing.T) {
	resp := Put("http://httpbin.org/get", nil)
	if resp != nil {
		t.Error("It is time to start writing the PUT tests")
	}

}
