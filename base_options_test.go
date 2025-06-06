package grequests

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestBasicOPTIONSRequest(t *testing.T) {
	resp, err := Options("http://httpbin.org/get")
	if err != nil {
		assert.Fail(t, "Unable to make OPTIONS request: ", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "Options request did not return success: ", resp.StatusCode)
	}

	if resp.Header.Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, PATCH, OPTIONS" {
		assert.Fail(t, "Access-Control-Allow-Methods Type Header is unexpected: ", resp.Header)
	}
}

func TestOptionsSession(t *testing.T) {
	session := NewSession(nil)

	resp, err := session.Options("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"one": "two"}})

	if err != nil {
		assert.FailNow(t, "Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	resp, err = session.Options("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"two": "three"}})

	if err != nil {
		assert.FailNow(t, "Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	resp, err = session.Options("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"three": "four"}})

	if err != nil {
		assert.FailNow(t, "Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	cookieURL, err := url.Parse("http://httpbin.org")
	if err != nil {
		assert.Fail(t, "We (for some reason) cannot parse the cookie URL")
	}

	for _, cookie := range session.HTTPClient.Jar.Cookies(cookieURL) {
		switch cookie.Name {
		case "one":
			if cookie.Value != "two" {
				assert.Fail(t, "Cookie value is not valid", cookie)
			}
		case "two":
			if cookie.Value != "three" {
				assert.Fail(t, "Cookie value is not valid", cookie)
			}
		case "three":
			if cookie.Value != "four" {
				assert.Fail(t, "Cookie value is not valid", cookie)
			}
		default:
			assert.Fail(t, "We should not have any other cookies: ", cookie)
		}
	}

}

func TestOptionsInvalidURLSession(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Options("%../dir/", nil); err == nil {
		assert.Fail(t, "Some how the request was valid to make request ", err)
	}
}
