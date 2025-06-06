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

func TestBasicPatchRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	// Example PATCH with some data
	ro := &RequestOptions{
		Data: map[string]string{"key": "value"}, // This will be sent as form data
	}
	resp, err := Patch(httpbinURL+"/patch", ro)

	if err != nil {
		assert.Fail(t, "Unable to make request", err, resp.Error)
	}

	assert.True(t, resp.Ok, "Request did not return OK. Status: ", resp.StatusCode, " Body: ", resp.String())

	// Verify response from local /patch endpoint
	var patchResp struct {
		Form map[string]string `json:"form"`
		URL  string            `json:"url"`
		// Add other fields like Headers, Origin, Data, JSON as needed for verification
	}
	err = resp.JSON(&patchResp)
	assert.NoError(t, err, "Failed to parse JSON response from /patch")
	assert.Equal(t, httpbinURL+"/patch", patchResp.URL)
	assert.Equal(t, "value", patchResp.Form["key"])
}

func TestPatchSession(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	session := NewSession(nil)

	// Set cookies
	_, err := session.Get(httpbinURL+"/cookies/set", &RequestOptions{Params: map[string]string{"one": "two"}})
	assert.NoError(t, err)
	_, err = session.Get(httpbinURL+"/cookies/set", &RequestOptions{Params: map[string]string{"two": "three"}})
	assert.NoError(t, err)
	_, err = session.Get(httpbinURL+"/cookies/set", &RequestOptions{Params: map[string]string{"three": "four"}})
	assert.NoError(t, err)

	// Make PATCH request with params (which will be sent as form data for PATCH)
	patchData := map[string]string{"patch_param": "patch_value"}
	ro := &RequestOptions{
		Data: patchData, // Using Data for PATCH body
	}
	patchResp, err := session.Patch(httpbinURL+"/patch", ro)

	assert.NoError(t, err, "PATCH request failed")
	assert.True(t, patchResp.Ok, "PATCH request did not return OK: ", patchResp.String())

	// Verify response from /patch
	var actualPatchRespContent struct {
		Form    map[string]string   `json:"form"` // httpbin_test_server returns map[string][]string for form
		URL     string              `json:"url"`
		Headers map[string][]string `json:"headers"`
	}
	// Adjusting expectation for Form to map[string][]string as per our server's behavior
	var actualPatchRespContentFixed struct {
		Form    map[string][]string `json:"form"`
		URL     string              `json:"url"`
		Headers map[string][]string `json:"headers"`
	}

	err = patchResp.JSON(&actualPatchRespContentFixed)
	assert.NoError(t, err, "Could not unmarshal PATCH response JSON: ", patchResp.String())
	assert.Equal(t, httpbinURL+"/patch", actualPatchRespContentFixed.URL)
	assert.Equal(t, []string{"patch_value"}, actualPatchRespContentFixed.Form["patch_param"])

	// Verify cookies are still present and sent with the PATCH request
	// The /patch endpoint will return the headers it received.
	// We need to check if the "Cookie" header was present in the request headers.
	cookieHeaderFound := false
	for key, values := range actualPatchRespContentFixed.Headers {
		if key == "Cookie" {
			for _, value := range values {
				assert.Contains(t, value, "one=two", "Cookie 'one' not found in PATCH request headers")
				assert.Contains(t, value, "two=three", "Cookie 'two' not found in PATCH request headers")
				assert.Contains(t, value, "three=four", "Cookie 'three' not found in PATCH request headers")
				cookieHeaderFound = true
			}
			break
		}
	}
	assert.True(t, cookieHeaderFound, "Cookie header not found in PATCH request to /patch")

	// Also check cookies in the session jar, they should persist
	parsedURL, err := url.Parse(httpbinURL)
	assert.NoError(t, err, "Cannot parse httpbinURL")
	cookiesFromJar := session.HTTPClient.Jar.Cookies(parsedURL)
	assert.Len(t, cookiesFromJar, 3, "Invalid number of cookies provided in Jar after PATCH")

	foundCookiesInJar := make(map[string]string)
	for _, cookie := range cookiesFromJar {
		foundCookiesInJar[cookie.Name] = cookie.Value
	}
	assert.Equal(t, "two", foundCookiesInJar["one"])
	assert.Equal(t, "three", foundCookiesInJar["two"])
	assert.Equal(t, "four", foundCookiesInJar["three"])
}

func TestPatchInvalidURLSession(t *testing.T) {
	// This test does not use httpbin, so it remains unchanged.
	session := NewSession(nil)

	if _, err := session.Patch("%../dir/", nil); err == nil {
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

func TestPatchInvalidURLSession(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Patch("%../dir/", nil); err == nil {
		assert.Fail(t, "Some how the request was valid to make request ", err)
	}
}
