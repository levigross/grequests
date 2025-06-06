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

func TestBasicPutRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	ro := &RequestOptions{
		Data: map[string]string{"key": "value"}, // Example data for PUT
	}
	resp, err := Put(httpbinURL+"/put", ro)

	if err != nil {
		assert.Fail(t, "Unable to make request", err, resp.Error)
	}
	assert.True(t, resp.Ok, "Request did not return OK. Status: ", resp.StatusCode, " Body: ", resp.String())

	// Verify response from local /put endpoint
	var putResp struct {
		Form map[string]string `json:"form"` // Assuming server returns form data like this
		URL  string            `json:"url"`
	}
	err = resp.JSON(&putResp)
	assert.NoError(t, err, "Failed to parse JSON response from /put")
	assert.Equal(t, httpbinURL+"/put", putResp.URL)
	assert.Equal(t, "value", putResp.Form["key"])
}

func TestBasicPutUploadRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	fileUploads, err := FileUploadFromDisk("testdata/mypassword")
	assert.NoError(t, err, "Unable to create FileUploadFromDisk")
	defer fileUploads[0].FileContents.Close()

	ro := &RequestOptions{
		Files: fileUploads,
		Data:  map[string]string{"One": "Two"},
	}
	resp, errPut := Put(httpbinURL+"/put", ro)

	assert.NoError(t, errPut, "Unable to make PUT request with file upload", resp.Error)
	assert.True(t, resp.Ok, "Request did not return OK. Status: ", resp.StatusCode, " Body: ", resp.String())

	var putResp struct {
		Files map[string]string      `json:"files"`
		Form  map[string]interface{} `json:"form"` // Using interface for flexibility
		URL   string                 `json:"url"`
	}
	err = resp.JSON(&putResp)
	assert.NoError(t, err, "Failed to parse JSON response from /put: ", resp.String())
	assert.Equal(t, httpbinURL+"/put", putResp.URL)
	assert.Equal(t, "saucy sauce", putResp.Files["mypassword"])

	formOne, ok := putResp.Form["One"].([]interface{})
	assert.True(t, ok, "Form field 'One' not found or not an array")
	if ok && len(formOne) > 0 {
		assert.Equal(t, "Two", formOne[0].(string))
	} else if ok {
		assert.Fail(t, "Form field 'One' was empty array")
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

	// Make PUT request with data
	putData := map[string]string{"put_key": "put_value"}
	ro := &RequestOptions{Data: putData}
	putResp, err := session.Put(httpbinURL+"/put", ro)

	assert.NoError(t, err, "PUT request failed")
	assert.True(t, putResp.Ok, "PUT request did not return OK: ", putResp.String())

	var actualPutRespContent struct {
		Form    map[string][]string `json:"form"` // Server returns form values as []string
		URL     string              `json:"url"`
		Headers map[string][]string `json:"headers"`
	}
	err = putResp.JSON(&actualPutRespContent)
	assert.NoError(t, err, "Could not unmarshal PUT response JSON: ", putResp.String())
	assert.Equal(t, httpbinURL+"/put", actualPutRespContent.URL)
	assert.Equal(t, []string{"put_value"}, actualPutRespContent.Form["put_key"])

	// Verify cookies were sent
	cookieHeaderFound := false
	for key, values := range actualPutRespContent.Headers {
		if key == "Cookie" {
			for _, value := range values {
				assert.Contains(t, value, "one=two")
				assert.Contains(t, value, "two=three")
				assert.Contains(t, value, "three=four")
				cookieHeaderFound = true
			}
			break
		}
	}
	assert.True(t, cookieHeaderFound, "Cookie header not found in PUT request to /put")

	// Check cookies in session jar
	parsedURL, err := url.Parse(httpbinURL)
	assert.NoError(t, err)
	cookiesFromJar := session.HTTPClient.Jar.Cookies(parsedURL)
	assert.Len(t, cookiesFromJar, 3, "Incorrect number of cookies in jar after PUT")

	foundCookiesInJar := make(map[string]string)
	for _, cookie := range cookiesFromJar {
		foundCookiesInJar[cookie.Name] = cookie.Value
	}
	assert.Equal(t, "two", foundCookiesInJar["one"])
	assert.Equal(t, "three", foundCookiesInJar["two"])
	assert.Equal(t, "four", foundCookiesInJar["three"])
}

func TestPutInvalidURLSession(t *testing.T) {
	// This test does not use httpbin, so it remains unchanged.
	session := NewSession(nil)

	if _, err := session.Put("%../dir/", nil); err == nil {
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
