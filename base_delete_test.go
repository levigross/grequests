package grequests

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

// setupHttpbinServerTest creates a new test server and returns its URL and a teardown function.
func setupHttpbinServerTest(t *testing.T) (string, func()) {
	ts := createHttpbinTestServer()
	return ts.URL, func() { ts.Close() }
}

func TestBasicDeleteRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, err := Delete(httpbinURL+"/delete", nil)

	if err != nil {
		assert.Fail(t, "Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	// Optional: verify response body if needed, similar to verifyOkResponse in base_get_test.go
	// For example:
	// var deleteResp struct {
	//	URL string `json:"url"`
	// }
	// err = resp.JSON(&deleteResp)
	// assert.NoError(t, err)
	// assert.Equal(t, httpbinURL+"/delete", deleteResp.URL)
}

func TestDeleteSession(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	session := NewSession(nil)

	resp, err := session.Get(httpbinURL+"/cookies/set", &RequestOptions{Params: map[string]string{"one": "two"}})
	assert.NoError(t, err, "Cannot set cookie 'one'")
	assert.True(t, resp.Ok, "Request to set cookie 'one' did not return OK")

	resp, err = session.Get(httpbinURL+"/cookies/set", &RequestOptions{Params: map[string]string{"two": "three"}})
	assert.NoError(t, err, "Cannot set cookie 'two'")
	assert.True(t, resp.Ok, "Request to set cookie 'two' did not return OK")

	resp, err = session.Get(httpbinURL+"/cookies/set", &RequestOptions{Params: map[string]string{"three": "four"}})
	assert.NoError(t, err, "Cannot set cookie 'three'")
	assert.True(t, resp.Ok, "Request to set cookie 'three' did not return OK")

	// Now make the DELETE request
	deleteResp, err := session.Delete(httpbinURL+"/delete", nil)
	assert.NoError(t, err, "Delete request failed")
	assert.True(t, deleteResp.Ok, "Delete request did not return OK")

	// Verify cookies are still present in the session jar
	parsedURL, err := url.Parse(httpbinURL)
	if err != nil {
		assert.FailNow(t, "We (for some reason) cannot parse the cookie URL: "+httpbinURL)
	}

	cookiesFromJar := session.HTTPClient.Jar.Cookies(parsedURL)
	assert.Len(t, cookiesFromJar, 3, "Invalid number of cookies provided in Jar")

	foundCookies := make(map[string]string)
	for _, cookie := range cookiesFromJar {
		foundCookies[cookie.Name] = cookie.Value
	}

	assert.Equal(t, "two", foundCookies["one"])
	assert.Equal(t, "three", foundCookies["two"])
	assert.Equal(t, "four", foundCookies["three"])

	// Verify the /delete response if necessary
	// var actualDeleteRespContent struct {
	//	Args    map[string]string `json:"args"`
	//	Headers http.Header       `json:"headers"`
	//	Origin  string            `json:"origin"`
	//	URL     string            `json:"url"`
	// }
	// err = deleteResp.JSON(&actualDeleteRespContent)
	// assert.NoError(t, err, "Could not unmarshal DELETE response JSON")
	// assert.Equal(t, httpbinURL+"/delete", actualDeleteRespContent.URL)
}

func TestDeleteInvalidURLSession(t *testing.T) {
	// This test does not use httpbin, so it remains unchanged.
	session := NewSession(nil)

	if _, err := session.Delete("%../dir/", nil); err == nil {
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
