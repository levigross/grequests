package grequests

import (
	"net/url"
	"testing"
)

func TestBasicOPTIONSRequest(t *testing.T) {
	resp, err := Options("http://httpbin.org/get", nil)
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

func TestOptionsSession(t *testing.T) {
	session := NewSession(nil)

	resp, err := session.Options("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"one": "two"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	resp, err = session.Options("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"two": "three"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	resp, err = session.Options("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"three": "four"}})

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

func TestOptionsInvalidURLSession(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Options("%../dir/", nil); err == nil {
		t.Error("Some how the request was valid to make request ", err)
	}
}
