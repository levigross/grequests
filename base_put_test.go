package grequests

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestBasicPutRequest(t *testing.T) {
	resp, err := Put("http://httpbin.org/put")

	if err != nil {
		assert.Fail(t, "Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

}

func TestBasicPutUploadRequest(t *testing.T) {
	fd, err := FileUploadFromDisk("testdata/mypassword")

	if err != nil {
		assert.Fail(t, "Unable to open file: ", err)
	}

	resp, _ := Put("http://httpbin.org/put",
		FromRequestOptions(&RequestOptions{
			Files: fd,
			Data:  map[string]string{"One": "Two"},
		}))

	if resp.Error != nil {
		assert.Fail(t, "Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

}

func TestBasicPutUploadRequestInvalidURL(t *testing.T) {
	fd, err := FileUploadFromDisk("testdata/mypassword")

	if err != nil {
		assert.Fail(t, "Unable to open file: ", err)
	}

	_, err = Put("%../dir/",
		FromRequestOptions(&RequestOptions{
			Files: fd,
			Data:  map[string]string{"One": "Two"},
		}))

	if err == nil {
		assert.FailNow(t, "Somehow able to make the request")
	}
}

func TestSessionPutUploadRequestInvalidURL(t *testing.T) {
	fd, err := FileUploadFromDisk("testdata/mypassword")

	if err != nil {
		assert.Fail(t, "Unable to open file: ", err)
	}

	session := NewSession(nil)

	_, err = session.Put("%../dir/",
		&RequestOptions{
			Files: fd,
			Data:  map[string]string{"One": "Two"},
		})

	if err == nil {
		assert.FailNow(t, "Somehow able to make the request")
	}
}

func TestPutSession(t *testing.T) {
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

	resp, err = session.Put("http://httpbin.org/put", &RequestOptions{Data: map[string]string{"one": "two"}})

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

func TestPutInvalidURLSession(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Put("%../dir/", nil); err == nil {
		assert.Fail(t, "Some how the request was valid to make request ", err)
	}
}
