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

func TestBasicHeadRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, err := Head(httpbinURL+"/get", nil) // Target /get, as HEAD should work for GET-able resources
	if err != nil {
		assert.Fail(t, "Unable to make HEAD request: ", err, resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "HEAD request did not return success: ", resp.StatusCode)
	}

	// http.DefaultServeMux typically sets Content-Type for HEAD if the GET handler would.
	// Our /get handler sets "application/json".
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json", "Content Type Header is unexpected")

	// Ensure no body for HEAD
	assert.Empty(t, resp.Bytes(), "HEAD request should not have a body (Bytes)")
	assert.Empty(t, resp.String(), "HEAD request should not have a body (String)")
}

func TestBasicHeadRequestNoContent(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	// Using /status/204 as it's designed for no content.
	// A HEAD request to /status/204 should also have no content and correct status.
	resp, err := Head(httpbinURL+"/status/204", nil)
	if err != nil {
		assert.Fail(t, "Unable to make HEAD request: ", err, resp.Error)
	}

	// httptest server + our /status/:code handler should result in 204.
	// For a HEAD request, it should still be 204.
	assert.Equal(t, 204, resp.StatusCode, "HEAD request to /status/204 did not return 204")
	assert.True(t, resp.Ok, "HEAD request did not return success classification for 204")


	// Content-Type for 204 responses is typically not set or empty.
	// Our /status/:code handler might set one based on fmt.Fprintf.
	// However, for 204, the body is empty.
	// Let's check what our server does. The server writes "204 No Content".
	// For a HEAD request, this body should be omitted.
	// Content-Length should be 0.
	assert.Equal(t, "0", resp.Header.Get("Content-Length"), "Content-Length should be 0 for HEAD to 204")


	if resp.Bytes() != nil && len(resp.Bytes()) > 0 { // Allow nil or empty byte slice
		assert.Fail(t, "Somehow byte buffer is working now (bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		assert.Fail(t, "Somehow byte buffer is working now (bytes)", resp.String())
	}
}

func TestHeadSession(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	session := NewSession(nil)

	// HEAD requests to /cookies/set will set cookies and follow redirect.
	// The final response will be from /cookies (after redirect).
	// A HEAD to /cookies should return headers of /cookies.
	resp, err := session.Head(httpbinURL+"/cookies/set", &RequestOptions{Params: map[string]string{"one": "two"}})
	assert.NoError(t, err, "Cannot set cookie 'one' via HEAD")
	// The response 'Ok' for a HEAD to a redirect depends on the final page after redirect.
	// /cookies returns 200 OK.
	assert.True(t, resp.Ok, "Request to set cookie 'one' did not return OK. Status: ", resp.StatusCode)


	resp, err = session.Head(httpbinURL+"/cookies/set", &RequestOptions{Params: map[string]string{"two": "three"}})
	assert.NoError(t, err, "Cannot set cookie 'two' via HEAD")
	assert.True(t, resp.Ok, "Request to set cookie 'two' did not return OK. Status: ", resp.StatusCode)


	resp, err = session.Head(httpbinURL+"/cookies/set", &RequestOptions{Params: map[string]string{"three": "four"}})
	assert.NoError(t, err, "Cannot set cookie 'three' via HEAD")
	assert.True(t, resp.Ok, "Request to set cookie 'three' did not return OK. Status: ", resp.StatusCode)


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
}

func TestHeadInvalidURLSession(t *testing.T) {
	// This test does not use httpbin, so it remains unchanged.
	session := NewSession(nil)

	if _, err := session.Head("%../dir/", nil); err == nil {
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
