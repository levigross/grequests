package grequests

import "testing"

func TestBasicReqRequestForGet (t *testing.T) {
	resp, _ := Req("GET","http://httpbin.org/get", nil)
	verifyOkResponse(resp, t)
}

func TestBasicReqRequestForDelete (t *testing.T)  {
	resp, err := Req("DELETE","http://httpbin.org/delete", nil)

	if err != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}
}

func TestBasicReqRequestForPost (t *testing.T) {
	resp, _ := Req("POST","http://httpbin.org/post",
		&RequestOptions{Data: map[string]string{"One": "Two"}})
	verifyOkPostResponse(resp, t)
}

func TestBasicReqRequestForPut (t *testing.T) {
	resp, err := Req("PUT","http://httpbin.org/put", nil)

	if err != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

}

func TestBasicReqRequestForPatch (t *testing.T) {
	resp, err := Req("PATCH","http://httpbin.org/patch", nil)

	if err != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

}

func TestBasicReqRequestForOptions (t *testing.T) {
	resp, err := Req("OPTIONS","http://httpbin.org/get", nil)
	if err != nil {
		t.Error("Unable to make OPTIONS request: ", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Options request did not return success: ", resp.StatusCode)
	}

	if resp.Header.Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, PATCH, OPTIONS" {
		t.Error("Access-Control-Allow-Methods Type Header is unexpected: ", resp.Header)
	}

}

func TestBasicReqRequestForHead (t *testing.T) {
	resp, err := Req("HEAD","http://httpbin.org/get", nil)
	if err != nil {
		t.Error("Unable to make HEAD request: ", resp.Error)
	}

	if resp.Ok != true {
		t.Error("HEAD request did not return success: ", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Error("Content Type Header is unexpected: ", resp.Header.Get("Content-Type"))
	}

}

