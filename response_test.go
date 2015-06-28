package grequests

import (
	"net/http"
	"strconv"
	"testing"
)

func TestResponseOk(t *testing.T) {
	status := []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}
	for _, status := range status {
		verifyResponseOkForStatus(status, t)
	}
}

func verifyResponseOkForStatus(status int, t *testing.T) {
	url := "http://httpbin.org/status/" + strconv.Itoa(status)
	resp, err := buildResponse(http.Get(url))

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Errorf("Request did not return OK. Received status code %d rather a 2xx.", resp.StatusCode)
	}
}
