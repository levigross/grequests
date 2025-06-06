package grequests

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestBasicHeadRequest(t *testing.T) {
	resp, err := Head("http://httpbin.org/get")
	if err != nil {
		assert.Fail(t, "Unable to make HEAD request: ", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "HEAD request did not return success: ", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		assert.Fail(t, "Content Type Header is unexpected: ", resp.Header.Get("Content-Type"))
	}
}

func TestBasicHeadRequestNoContent(t *testing.T) {
	resp, err := Head("http://httpbin.org/bytes/0")
	if err != nil {
		assert.Fail(t, "Unable to make HEAD request: ", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "HEAD request did not return success: ", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/octet-stream" {
		assert.Fail(t, "Content Type Header is unexpected: ", resp.Header.Get("Content-Type"))
	}

	if resp.Bytes() != nil {
		assert.Fail(t, "Somehow byte buffer is working now (bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		assert.Fail(t, "Somehow byte buffer is working now (bytes)", resp.String())
	}
}

func TestHeadSession(t *testing.T) {
	session := NewSession(nil)

	resp, err := session.Head("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"one": "two"}})

	if err != nil {
		assert.FailNow(t, "Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	resp, err = session.Head("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"two": "three"}})

	if err != nil {
		assert.FailNow(t, "Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	resp, err = session.Head("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"three": "four"}})

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

	if len(session.HTTPClient.Jar.Cookies(cookieURL)) != 3 {
		assert.Fail(t, "Invalid number of cookies provided: ", session.HTTPClient.Jar.Cookies(cookieURL))
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

func TestHeadInvalidURLSession(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Head("%../dir/", nil); err == nil {
		assert.Fail(t, "Some how the request was valid to make request ", err)
	}
}
