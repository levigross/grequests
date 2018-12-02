package grequests

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestBasicDeleteRequest(t *testing.T) {
	resp, err := Delete("http://httpbin.org/delete")

	if err != nil {
		assert.Fail(t, "Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}
}

func TestDeleteSession(t *testing.T) {
	session := NewSession(nil)

	resp, err := session.Get("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"one": "two"}})

	if err != nil {
		assert.FailNow(t, "Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	resp, err = session.Get("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"two": "three"}})

	if err != nil {
		assert.FailNow(t, "Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	resp, err = session.Get("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"three": "four"}})

	if err != nil {
		assert.FailNow(t, "Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	resp, err = session.Delete("http://httpbin.org/delete", nil)

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
		assert.Fail(t, "Invalid number of cookies provided: ", resp.RawResponse.Cookies())
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

func TestDeleteInvalidURLSession(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Delete("%../dir/", nil); err == nil {
		assert.Fail(t, "Some how the request was valid to make request ", err)
	}
}
