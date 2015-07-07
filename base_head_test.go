package grequests

import (
	"net/url"
	"testing"
)

func TestBasicHeadRequest(t *testing.T) {
	resp, err := Head("http://httpbin.org/get", nil)
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

func TestBasicHeadRequestNoContent(t *testing.T) {
	resp, err := Head("http://httpbin.org/bytes/0", nil)
	if err != nil {
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

func TestHeadSession(t *testing.T) {
	session := NewSession(nil)

	resp, err := session.Head("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"one": "two"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	resp, err = session.Head("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"two": "three"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	resp, err = session.Head("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"three": "four"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	cookieURL, err := url.Parse("http://httpbin.org")
	if err != nil {
		t.Error("We (for some reason) cannot parse the cookie URL")
	}

	if len(session.HTTPClient.Jar.Cookies(cookieURL)) != 3 {
		t.Error("Invalid number of cookies provided: ", session.HTTPClient.Jar.Cookies(cookieURL))
	}

	for _, cookie := range session.HTTPClient.Jar.Cookies(cookieURL) {
		switch cookie.Name {
		case "one":
			if cookie.Value != "two" {
				t.Error("Cookie value is not valid", cookie)
			}
		case "two":
			if cookie.Value != "three" {
				t.Error("Cookie value is not valid", cookie)
			}
		case "three":
			if cookie.Value != "four" {
				t.Error("Cookie value is not valid", cookie)
			}
		default:
			t.Error("We should not have any other cookies: ", cookie)
		}
	}

}

func TestHeadInvalidURLSession(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Head("%../dir/", nil); err == nil {
		t.Error("Some how the request was valid to make request ", err)
	}
}
