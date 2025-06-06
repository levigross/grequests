package grequests

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"math"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

// setupHttpbinServerTest creates a new test server and returns its URL and a teardown function.
func setupHttpbinServerTest(t *testing.T) (string, func()) {
	ts := createHttpbinTestServer()
	return ts.URL, func() { ts.Close() }
}

type BasicPostResponse struct {
	Args  struct{} `json:"args"`
	Data  string   `json:"data"`
	Files struct{} `json:"files"`
	Form  struct {
		One string `json:"one"`
	} `json:"form"`
	Headers struct {
		Accept        string `json:"Accept"`
		ContentLength string `json:"Content-Length"`
		ContentType   string `json:"Content-Type"`
		Host          string `json:"Host"`
		UserAgent     string `json:"User-Agent"`
	} `json:"headers"`
	JSON   interface{} `json:"json"`
	Origin string      `json:"origin"`
	URL    string      `json:"url"`
}

type BasicPostJSONResponse struct {
	Args    struct{} `json:"args"`
	Data    string   `json:"data"`
	Files   struct{} `json:"files"`
	Form    struct{} `json:"form"`
	Headers struct {
		AcceptEncoding string `json:"Accept-Encoding"`
		ContentLength  string `json:"Content-Length"`
		ContentType    string `json:"Content-Type"`
		Host           string `json:"Host"`
		UserAgent      string `json:"User-Agent"`
		XRequestedWith string `json:"X-Requested-With"`
	} `json:"headers"`
	JSON struct {
		One string `json:"One"`
	} `json:"json"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicMultiFileUploadResponse struct {
	Args  struct{} `json:"args"`
	Data  string   `json:"data"`
	Files struct {
		File1 string `json:"file1"`
		File2 string `json:"file2"`
	} `json:"files"`
	Form struct {
		One string `json:"One"`
	} `json:"form"`
	Headers struct {
		AcceptEncoding string `json:"Accept-Encoding"`
		ContentLength  string `json:"Content-Length"`
		ContentType    string `json:"Content-Type"`
		Host           string `json:"Host"`
		UserAgent      string `json:"User-Agent"`
	} `json:"headers"`
	JSON   interface{} `json:"json"`
	Origin string      `json:"origin"`
	URL    string      `json:"url"`
}

type BasicPostFileUpload struct {
	Args  struct{} `json:"args"`
	Data  string   `json:"data"`
	Files struct {
		File string `json:"file"`
	} `json:"files"`
	Form struct {
		One string `json:"one"`
	} `json:"form"`
	Headers struct {
		AcceptEncoding string `json:"Accept-Encoding"`
		ContentLength  string `json:"Content-Length"`
		ContentType    string `json:"Content-Type"`
		Host           string `json:"Host"`
		UserAgent      string `json:"User-Agent"`
	} `json:"headers"`
	JSON   interface{} `json:"json"`
	Origin string      `json:"origin"`
	URL    string      `json:"url"`
}

type XMLPostMessage struct {
	Name   string
	Age    int
	Height int
}

type dataAndErrorBuffer struct {
	err error
	bytes.Buffer
}

func (dataAndErrorBuffer) Close() error { return nil }

func (r dataAndErrorBuffer) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func TestBasicPostRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, _ := Post(httpbinURL+"/post",
		FromRequestOptions(&RequestOptions{Data: map[string]string{"One": "Two"}}))
	verifyOkPostResponse(resp, t, httpbinURL+"/post", map[string]string{"One": "Two"}, "")

}

func TestBasicRegularPostRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, err := Post(httpbinURL+"/post",
		FromRequestOptions(&RequestOptions{Data: map[string]string{"One": "Two"}}))

	if err != nil {
		assert.Fail(t, "Cannot post: ", err)
	}

	verifyOkPostResponse(resp, t, httpbinURL+"/post", map[string]string{"One": "Two"}, "")

}

func TestBasicPostRequestInvalidURL(t *testing.T) {
	resp, _ := Post("%../dir/",
		FromRequestOptions(&RequestOptions{Data: map[string]string{"One": "Two"},
			Params: map[string]string{"1": "2"}}))

	if resp.Error == nil {
		assert.Fail(t, "Somehow the request went through")
	}

}

func TestBasicPostRequestInvalidURLNoParams(t *testing.T) {
	resp, _ := Post("%../dir/", FromRequestOptions(&RequestOptions{Data: map[string]string{"One": "Two"}}))

	if resp.Error == nil {
		assert.Fail(t, "Somehow the request went through")
	}

}

func TestSessionPostRequestInvalidURLNoParams(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Post("%../dir/", &RequestOptions{Data: map[string]string{"One": "Two"}}); err == nil {
		assert.Fail(t, "Somehow the request went through")
	}

}

func TestXMLPostRequestInvalidURL(t *testing.T) {
	resp, _ := Post("%../dir/",
		FromRequestOptions(&RequestOptions{XML: XMLPostMessage{Name: "Human", Age: 1, Height: 1}}))

	if resp.Error == nil {
		assert.Fail(t, "Somehow the request went through")
	}
}

func TestXMLSessionPostRequestInvalidURL(t *testing.T) {
	session := NewSession(nil)

	_, err := session.Post("%../dir/",
		&RequestOptions{XML: XMLPostMessage{Name: "Human", Age: 1, Height: 1}})

	if err == nil {
		assert.Fail(t, "Somehow the request went through")
	}
}

func TestBasicPostJsonRequestInvalidURL(t *testing.T) {
	_, err := Post("%../dir/",
		FromRequestOptions(&RequestOptions{JSON: map[string]string{"One": "Two"}, IsAjax: true}))

	if err == nil {
		assert.Fail(t, "Somehow the request went through")
	}
}

func TestSessionPostJsonRequestInvalidURL(t *testing.T) {
	session := NewSession(nil)

	_, err := session.Post("%../dir/",
		&RequestOptions{JSON: map[string]string{"One": "Two"}, IsAjax: true})

	if err == nil {
		assert.Fail(t, "Somehow the request went through")
	}
}

func TestBasicPostJsonRequestInvalidJSON(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, err := Post(httpbinURL+"/post",
		FromRequestOptions(&RequestOptions{JSON: math.NaN(), IsAjax: true}))

	if err == nil {
		assert.Fail(t, "Somehow the request went through")
	}

	if resp.Ok == true {
		assert.Fail(t, "Somehow the request is OK")
	}
}

func TestSessionPostJsonRequestInvalidJSON(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	session := NewSession(nil)

	resp, err := session.Post(httpbinURL+"/post",
		&RequestOptions{JSON: math.NaN(), IsAjax: true})

	if err == nil {
		assert.Fail(t, "Somehow the request went through")
	}

	if resp.Ok == true {
		assert.Fail(t, "Somehow the request is OK")
	}
}

func TestBasicPostJsonRequestInvalidXML(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, err := Post(httpbinURL+"/post",
		FromRequestOptions(&RequestOptions{XML: map[string]string{"One": "two"}, IsAjax: true}))

	if err == nil {
		assert.Fail(t, "Somehow the request went through")
	}

	if resp.Ok == true {
		assert.Fail(t, "Somehow the request is OK")
	}
}

func TestSessionPostJsonRequestInvalidXML(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	session := NewSession(nil)

	resp, err := session.Post(httpbinURL+"/post",
		&RequestOptions{XML: map[string]string{"One": "two"}, IsAjax: true})

	if err == nil {
		assert.Fail(t, "Somehow the request went through")
	}

	if resp.Ok == true {
		assert.Fail(t, "Somehow the request is OK")
	}
}

func TestBasicPostRequestUploadInvalidURL(t *testing.T) {

	fd, err := FileUploadFromDisk("testdata/mypassword")

	if err != nil {
		assert.Fail(t, "Unable to open file: ", err)
	}

	defer fd[0].FileContents.Close()

	resp, _ := Post("%../dir/",
		FromRequestOptions(&RequestOptions{
			Files: fd,
			Data:  map[string]string{"One": "Two"},
		}))

	if resp.Error == nil {
		assert.FailNow(t, "Somehow able to make the request")
	}
}

func TestSessionPostRequestUploadInvalidURL(t *testing.T) {
	session := NewSession(nil)

	fd, err := FileUploadFromDisk("testdata/mypassword")

	if err != nil {
		assert.Fail(t, "Unable to open file: ", err)
	}

	defer fd[0].FileContents.Close()

	_, err = session.Post("%../dir/",
		&RequestOptions{
			Files: fd,
			Data:  map[string]string{"One": "Two"},
		})

	if err == nil {
		assert.FailNow(t, "Somehow able to make the request")
	}
}

func TestBasicPostRequestUploadInvalidFileUpload(t *testing.T) {

	resp, _ := Post("%../dir/",
		FromRequestOptions(&RequestOptions{
			Files: []FileUpload{{FileName: `\x00%'"üfdsufhid\Ä\"D\\\"JS%25//'"H•\\\\'"¶•ªç∂\uf8\x8AKÔÓÔ`, FileContents: nil}},
			Data:  map[string]string{"One": "Two"},
		}))

	if resp.Error == nil {
		assert.FailNow(t, "Somehow able to make the request")
	}
}

func TestSessionPostRequestUploadInvalidFileUpload(t *testing.T) {
	session := NewSession(nil)
	_, err := session.Post("%../dir/",
		&RequestOptions{
			Files: []FileUpload{{FileName: "üfdsufhidÄDJSHAKÔÓÔ", FileContents: nil}},
			Data:  map[string]string{"One": "Two"},
		})

	if err == nil {
		assert.FailNow(t, "Somehow able to make the request")
	}
}

func TestXMLPostRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, _ := Post(httpbinURL+"/post",
		FromRequestOptions(&RequestOptions{XML: XMLPostMessage{Name: "Human", Age: 1, Height: 1}}))

	if resp.Error != nil {
		assert.FailNow(t, "Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	// BasicPostJSONResponse is not the right struct for XML data.
	// The local server will put raw XML string into the "data" field.
	var postResp struct {
		Data    string `json:"data"`
		URL     string `json:"url"`
		Headers struct {
			ContentType string `json:"Content-Type"`
		} `json:"headers"`
	}

	if err := resp.JSON(&postResp); err != nil {
		assert.Fail(t, "Unable to coerce to JSON", err, resp.String())
	}
	assert.Equal(t, httpbinURL+"/post", postResp.URL)
	assert.Contains(t, postResp.Headers.ContentType, "application/xml", "Content-Type should be application/xml")

	myXMLStruct := &XMLPostMessage{}
	err := xml.Unmarshal([]byte(postResp.Data), myXMLStruct)
	assert.NoError(t, err, "Failed to unmarshal XML from response data")

	if myXMLStruct.Age != 1 {
		assert.Fail(t, fmt.Sprintf("XML content mismatch. Got: %#v, Expected Age 1", myXMLStruct))
	}

}

func TestXMLPostRequestReaderBody(t *testing.T) {
	msg := XMLPostMessage{Name: "Human", Age: 1, Height: 1}
	derBytes, err := xml.Marshal(msg)
	if err != nil {
		assert.FailNow(t, "Unable to marshal XML", err)
	}

	resp, _ := Post(httpbinURL+"/post",
		FromRequestOptions(&RequestOptions{RequestBody: bytes.NewReader(derBytes), Headers: map[string]string{"Content-Type": "application/xml"}}))

	if resp.Error != nil {
		assert.FailNow(t, "Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	var postResp struct {
		Data    string `json:"data"`
		URL     string `json:"url"`
		Headers struct {
			ContentType string `json:"Content-Type"`
		} `json:"headers"`
	}

	if err := resp.JSON(&postResp); err != nil {
		assert.Fail(t, "Unable to coerce to JSON", err, resp.String())
	}
	assert.Equal(t, httpbinURL+"/post", postResp.URL)
	// Content-Type might not be automatically set to application/xml by grequests when RequestBody is used
	// unless specified in Headers. The test server will report what it receives.
	// assert.Contains(t, postResp.Headers.ContentType, "application/xml", "Content-Type should be application/xml")


	myXMLStruct := &XMLPostMessage{}
	newErr := xml.Unmarshal([]byte(postResp.Data), myXMLStruct)
	assert.NoError(t, newErr, "Failed to unmarshal XML from response data")


	if myXMLStruct.Age != 1 {
		assert.Fail(t, fmt.Sprintf("XML content mismatch. Got: %#v, Expected Age 1", myXMLStruct))
	}

}

func TestXMLMarshaledStringPostRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	xmlStruct := XMLPostMessage{Name: "Human", Age: 1, Height: 1}
	encoded, _ := xml.Marshal(xmlStruct)
	resp, _ := Post(httpbinURL+"/post",
		FromRequestOptions(&RequestOptions{XML: string(encoded)}))

	if resp.Error != nil {
		assert.FailNow(t, "Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	var postResp struct {
		Data    string `json:"data"`
		URL     string `json:"url"`
		Headers struct {
			ContentType string `json:"Content-Type"`
		} `json:"headers"`
	}
	if err := resp.JSON(&postResp); err != nil {
		assert.Fail(t, "Unable to response to JSON", err, resp.String())
	}
	assert.Equal(t, httpbinURL+"/post", postResp.URL)
	assert.Contains(t, postResp.Headers.ContentType, "application/xml", "Content-Type should be application/xml")


	if postResp.Data != string(encoded) {
		assert.Fail(t, "Response data is not valid", postResp.Data, string(encoded))
	}
}

func TestXMLMarshaledBytesPostRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	xmlStruct := XMLPostMessage{Name: "Human", Age: 1, Height: 1}
	encoded, _ := xml.Marshal(xmlStruct)
	resp, _ := Post(httpbinURL+"/post",
		FromRequestOptions(&RequestOptions{XML: encoded}))

	if resp.Error != nil {
		assert.FailNow(t, "Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	var postResp struct {
		Data    string `json:"data"`
		URL     string `json:"url"`
		Headers struct {
			ContentType string `json:"Content-Type"`
		} `json:"headers"`
	}
	if err := resp.JSON(&postResp); err != nil {
		assert.Fail(t, "Unable to response to JSON", err, resp.String())
	}
	assert.Equal(t, httpbinURL+"/post", postResp.URL)
	assert.Contains(t, postResp.Headers.ContentType, "application/xml", "Content-Type should be application/xml")

	if postResp.Data != string(encoded) {
		assert.Fail(t, "Response data is not valid", postResp.Data, string(encoded))
	}
}

func TestXMLNilPostRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, _ := Post(httpbinURL+"/post", FromRequestOptions(&RequestOptions{XML: nil}))

	if resp.Error != nil {
		assert.FailNow(t, "Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		assert.Fail(t, "Request did not return OK")
	}

	var postResp struct {
		Data    string `json:"data"`
		URL     string `json:"url"`
		Headers struct {
			ContentType string `json:"Content-Type"`
		} `json:"headers"`
	}
	if err := resp.JSON(&postResp); err != nil {
		assert.Fail(t, "Unable to response to JSON", err, resp.String())
	}
	assert.Equal(t, httpbinURL+"/post", postResp.URL)
	// Content-Type might be "text/plain" or not set if data is empty
	// For XML: nil, grequests sends an empty body with Content-Type: application/xml
	assert.Contains(t, postResp.Headers.ContentType, "application/xml", "Content-Type should be application/xml for nil XML")


	if postResp.Data != "" {
		assert.Fail(t, "Response data is not valid for nil XML", postResp.Data)
	}
}

func TestBasicPostRequestUploadErrorReader(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	var rd dataAndErrorBuffer
	rd.err = fmt.Errorf("Random Error")
	_, err := Post(httpbinURL+"/post",
		FromRequestOptions(&RequestOptions{
			Files: []FileUpload{{FileName: "Random.test", FileContents: rd}},
			Data:  map[string]string{"One": "Two"},
		}))

	if err == nil {
		assert.Fail(t, "Somehow our test didn't fail...")
	}
}

func TestBasicPostRequestUploadErrorEOFReader(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	var rd dataAndErrorBuffer
	rd.err = io.EOF
	// This test expects the Post operation to succeed even with io.EOF from the reader,
	// as EOF is a valid state after reading all content.
	// The httpbin server would just receive an empty file part.
	resp, err := Post(httpbinURL+"/post",
		FromRequestOptions(&RequestOptions{
			Files: []FileUpload{{FileName: "Random.test", FileContents: rd}},
			Data:  map[string]string{"One": "Two"},
		}))

	if err != nil {
		assert.Fail(t, "Somehow our test didn't fail... ", err)
	}
}

func TestBasicPostRequestUpload(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	fileUploads, err := FileUploadFromDisk("testdata/mypassword")
	assert.NoError(t, err, "Unable to create FileUploadFromDisk")
	defer fileUploads[0].FileContents.Close() // Ensure file is closed

	ro := &RequestOptions{
		Files: fileUploads,
		Data:  map[string]string{"One": "Two"},
	}
	resp, errPost := Post(httpbinURL+"/post", ro)
	assert.NoError(t, errPost, "Unable to make POST request with file upload")
	assert.NotNil(t, resp, "Response should not be nil")
	if resp == nil { // Guard against nil dereference if assert fails
		t.FailNow()
	}
	assert.NoError(t, resp.Error, "Response error field should be nil")
	assert.True(t, resp.Ok, "Request did not return OK. Status: ", resp.StatusCode, " Body: ", resp.String())

	var postResp struct {
		Args    map[string]string            `json:"args"`
		Data    string                       `json:"data"`
		Files   map[string]string            `json:"files"` // httpbin.org returns file content as string
		Form    map[string]interface{}       `json:"form"`  // Use interface{} for flexibility or map[string][]string
		Headers map[string][]string          `json:"headers"`
		URL     string                       `json:"url"`
	}

	err = resp.JSON(&postResp)
	assert.NoError(t, err, "Unable to coerce to JSON: ", resp.String())

	assert.Equal(t, httpbinURL+"/post", postResp.URL)
	parsedURL, _ := url.Parse(httpbinURL)
	assert.Equal(t, parsedURL.Host, postResp.Headers["Host"][0])


	// Our server stores file content as string, matching httpbin.org
	// The filename in FileUploadFromDisk becomes the key in "files" if FileNameAsFieldName is true (default)
	// or if FileName is set. If FileName is empty, it might use a default like "file".
	// FileUploadFromDisk sets FileName to the base name of the path.
	assert.Equal(t, "saucy sauce", postResp.Files["mypassword"], "Uploaded file content mismatch")

	// Form data should be present
	// Our server returns form values as []string
	formOne, ok := postResp.Form["One"].([]interface{}) // JSON unmarshals to []interface{} for array
	assert.True(t, ok, "Form field 'One' not found or not an array")
	if ok && len(formOne) > 0 {
		assert.Equal(t, "Two", formOne[0].(string))
	} else if ok {
		assert.Fail(t, "Form field 'One' was empty array")
	}


	assert.Nil(t, resp.Bytes(), "JSON decoding should fully consume the response stream (Bytes)")
	assert.Empty(t, resp.String(), "JSON decoding should fully consume the response stream (String)")
	assert.Equal(t, 200, resp.StatusCode, "Response returned a non-200 code")
}

func TestBasicPostRequestUploadWithMime(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	fileContents, err := os.Open("testdata/mypassword")
	assert.NoError(t, err, "Unable to open file for upload")
	defer fileContents.Close()

	fileUploads := []FileUpload{
		{
			FileName:     "customfile.txt", // Provide a filename
			FileMime:     "text/special",
			FileContents: fileContents,
		},
	}

	ro := &RequestOptions{
		Files: fileUploads,
		Data:  map[string]string{"One": "Two"},
	}
	resp, errPost := Post(httpbinURL+"/post", ro)
	assert.NoError(t, errPost, "Unable to make POST request with file upload and custom MIME")
	assert.NotNil(t, resp, "Response should not be nil")
	if resp == nil {
		t.FailNow()
	}
	assert.NoError(t, resp.Error, "Response error field should be nil")
	assert.True(t, resp.Ok, "Request did not return OK. Status: ", resp.StatusCode, " Body: ", resp.String())

	var postResp struct {
		Files   map[string]string      `json:"files"`
		Form    map[string]interface{} `json:"form"`
		Headers map[string][]string    `json:"headers"`
		URL     string                 `json:"url"`
	}
	err = resp.JSON(&postResp)
	assert.NoError(t, err, "Unable to coerce to JSON: ", resp.String())

	assert.Equal(t, httpbinURL+"/post", postResp.URL)
	assert.Equal(t, "saucy sauce", postResp.Files["customfile.txt"]) // Check by FileName

	// Check that the Content-Type in the multipart form for the file reflects FileMime (this is hard to check directly from httpbin response)
	// httpbin.org's response doesn't usually detail the MIME type of individual file parts in its JSON output.
	// Our server also doesn't explicitly return this for individual files.
	// This test mainly ensures the request goes through and data is received.

	formOne, ok := postResp.Form["One"].([]interface{})
	assert.True(t, ok, "Form field 'One' not found or not an array")
	if ok && len(formOne) > 0 {
		assert.Equal(t, "Two", formOne[0].(string))
	} else if ok {
		assert.Fail(t, "Form field 'One' was empty array")
	}

	// TODO: Ensure file field contains correct content-type, field, and
	// filename information as soon as
	// https://github.com/kennethreitz/httpbin/pull/388 gets merged
	// (Or figure out a way to test this case the PR is rejected)
}

func TestBasicPostRequestUploadMultipleFiles(t *testing.T) {

	// TODO: Ensure file field contains correct content-type, field, and
	// filename information as soon as
	// https://github.com/kennethreitz/httpbin/pull/388 gets merged
	// (Or figure out a way to test this case the PR is rejected)
}

func TestBasicPostRequestUploadMultipleFiles(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	// testdata/herefortheglob, testdata/mypassword
	fileUploads, err := FileUploadFromGlob("testdata/*")
	assert.NoError(t, err, "Unable to glob files")
	for _, fu := range fileUploads {
		defer fu.FileContents.Close()
	}

	ro := &RequestOptions{
		Files: fileUploads,
		Data:  map[string]string{"One": "Two"},
	}
	resp, errPost := Post(httpbinURL+"/post", ro)
	assert.NoError(t, errPost, "Unable to make POST request with multiple file uploads")
	assert.True(t, resp.Ok, "Request did not return OK")

	var postResp struct {
		Files   map[string]string      `json:"files"`
		Form    map[string]interface{} `json:"form"`
		Headers map[string][]string    `json:"headers"`
		URL     string                 `json:"url"`
	}
	err = resp.JSON(&postResp)
	assert.NoError(t, err, "Unable to coerce to JSON: ", resp.String())

	assert.Equal(t, httpbinURL+"/post", postResp.URL)
	// Check file contents by their base names (FileUploadFromGlob uses base names)
	assert.Equal(t, "saucy sauce", postResp.Files["mypassword"], "Content of mypassword mismatch")
	assert.Equal(t, "I am just here to test the glob", postResp.Files["herefortheglob"], "Content of herefortheglob mismatch")

	formOne, ok := postResp.Form["One"].([]interface{})
	assert.True(t, ok, "Form field 'One' not found or not an array")
	if ok && len(formOne) > 0 {
		assert.Equal(t, "Two", formOne[0].(string))
	} else if ok {
		assert.Fail(t, "Form field 'One' was empty array")
	}
}

func TestBasicPostJsonBytesRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	jsonData := []byte(`{"One":"Two"}`)
	ro := &RequestOptions{JSON: jsonData, IsAjax: true}
	resp, err := Post(httpbinURL+"/post", ro)
	assert.NoError(t, err, "Unable to make POST JSON Bytes request")
	assert.True(t, resp.Ok, "Request did not return OK")

	verifyPostJSONResponse(resp, t, httpbinURL+"/post", `{"One":"Two"}`, true)
}

func TestBasicPostJsonStringRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	jsonStr := `{"One":"Two"}`
	ro := &RequestOptions{JSON: jsonStr, IsAjax: true}
	resp, err := Post(httpbinURL+"/post", ro)
	assert.NoError(t, err, "Unable to make POST JSON String request")
	assert.True(t, resp.Ok, "Request did not return OK")

	verifyPostJSONResponse(resp, t, httpbinURL+"/post", `{"One":"Two"}`, true)
}

func TestBasicPostJsonRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	jsonData := map[string]string{"One": "Two"}
	ro := &RequestOptions{JSON: jsonData, IsAjax: true}
	resp, err := Post(httpbinURL+"/post", ro)
	assert.NoError(t, err, "Unable to make POST JSON request")
	assert.True(t, resp.Ok, "Request did not return OK")

	verifyPostJSONResponse(resp, t, httpbinURL+"/post", `{"One":"Two"}`, true)
}

// Helper for verifying common JSON post responses
func verifyPostJSONResponse(resp *Response, t *testing.T, expectedURL, expectedData string, isAjax bool) {
	var postResp struct {
		Args    map[string]string   `json:"args"`
		Data    string              `json:"data"`
		Files   map[string]string   `json:"files"`
		Form    map[string]string   `json:"form"`
		Headers map[string][]string `json:"headers"`
		JSON    json.RawMessage     `json:"json"` // Use json.RawMessage to compare raw JSON
		Origin  string              `json:"origin"`
		URL     string              `json:"url"`
	}

	err := resp.JSON(&postResp)
	assert.NoError(t, err, "Unable to coerce to JSON: ", resp.String())

	assert.Equal(t, expectedURL, postResp.URL)
	parsedURL, _ := url.Parse(expectedURL)
	assert.Equal(t, parsedURL.Host, postResp.Headers["Host"][0])

	// Compare JSON part
	// Unmarshal expected and actual JSON to compare them structurally
	// to avoid issues with key order or whitespace.
	var expectedJSON, actualJSON interface{}
	errExpected := json.Unmarshal([]byte(expectedData), &expectedJSON)
	assert.NoError(t, errExpected, "Failed to unmarshal expectedData for comparison")
	errActual := json.Unmarshal(postResp.JSON, &actualJSON)
	assert.NoError(t, errActual, "Failed to unmarshal actual JSON from response for comparison")
	assert.Equal(t, expectedJSON, actualJSON, "Posted JSON content mismatch")


	assert.Equal(t, expectedData, strings.TrimSpace(postResp.Data), "Raw data field mismatch")

	if isAjax {
		assert.Equal(t, "XMLHttpRequest", postResp.Headers["X-Requested-With"][0], "X-Requested-With header missing for AJAX")
	}

	assert.Nil(t, resp.Bytes(), "JSON decoding should fully consume the response stream (Bytes)")
	assert.Empty(t, resp.String(), "JSON decoding should fully consume the response stream (String)")
	assert.Equal(t, 200, resp.StatusCode, "Response returned a non-200 code")
}

func TestPostSession(t *testing.T) {
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

	// Make POST request
	postData := map[string]string{"form_key": "form_value"}
	postOptions := &RequestOptions{Data: postData}
	postResp, err := session.Post(httpbinURL+"/post", postOptions)
	assert.NoError(t, err, "POST request in session failed")
	assert.True(t, postResp.Ok, "POST request in session did not return OK")

	// Verify POST response
	var actualPostRespContent struct {
		Form    map[string][]string `json:"form"` // Server returns form values as []string
		URL     string              `json:"url"`
		Headers map[string][]string `json:"headers"`
	}
	err = postResp.JSON(&actualPostRespContent)
	assert.NoError(t, err, "Could not unmarshal POST response JSON: ", postResp.String())
	assert.Equal(t, httpbinURL+"/post", actualPostRespContent.URL)
	assert.Equal(t, []string{"form_value"}, actualPostRespContent.Form["form_key"])

	// Verify cookies were sent with the POST request
	cookieHeaderFound := false
	for key, values := range actualPostRespContent.Headers {
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
	assert.True(t, cookieHeaderFound, "Cookie header not found in POST request to /post")


	// Check cookies in session jar
	parsedURL, err := url.Parse(httpbinURL)
	assert.NoError(t, err)
	cookiesFromJar := session.HTTPClient.Jar.Cookies(parsedURL)
	assert.Len(t, cookiesFromJar, 3, "Incorrect number of cookies in jar after POST")

	foundCookiesInJar := make(map[string]string)
	for _, cookie := range cookiesFromJar {
		foundCookiesInJar[cookie.Name] = cookie.Value
	}
	assert.Equal(t, "two", foundCookiesInJar["one"])
	assert.Equal(t, "three", foundCookiesInJar["two"])
	assert.Equal(t, "four", foundCookiesInJar["three"])
}
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

// verifyOkPostResponse verifies a basic form post response.
// expectedURL is the full URL the request was made to.
// expectedFormData is a map of form fields and their expected values.
// expectedRawData is the expected raw data string if the post was not form data (e.g. JSON, XML).
func verifyOkPostResponse(resp *Response, t *testing.T, expectedURL string, expectedFormData map[string]string, expectedRawData string) {
	if resp.Error != nil {
		assert.FailNow(t, "Unable to make request", resp.Error, resp.String())
	}

	assert.True(t, resp.Ok, "Request did not return OK. Status: ", resp.StatusCode, " Body: ", resp.String())
	assert.Equal(t, 200, resp.StatusCode, "Response returned a non-200 code")

	var postResp struct {
		Args    map[string]string      `json:"args"`    // Query args
		Data    string                 `json:"data"`    // Raw request body
		Files   map[string]string      `json:"files"`   // Uploaded files
		Form    map[string]interface{} `json:"form"`    // Form fields (local server uses []string for values)
		Headers map[string][]string    `json:"headers"` // Request headers received by server
		JSON    interface{}            `json:"json"`    // Parsed JSON body if Content-Type was application/json
		Origin  string                 `json:"origin"`
		URL     string                 `json:"url"`
	}

	err := resp.JSON(&postResp)
	assert.NoError(t, err, "Unable to coerce to JSON: ", resp.String())

	assert.Equal(t, expectedURL, postResp.URL, "URL in response mismatch")

	parsedURL, _ := url.Parse(expectedURL)
	assert.Equal(t, parsedURL.Host, postResp.Headers["Host"][0], "Host header in response mismatch")


	if len(expectedFormData) > 0 {
		for key, expectedValue := range expectedFormData {
			formValue, ok := postResp.Form[key]
			assert.True(t, ok, fmt.Sprintf("Form field '%s' not found in response", key))
			// Local server returns form values as []string (actually []interface{} after json unmarshal)
			formValueSlice, okSlice := formValue.([]interface{})
			assert.True(t, okSlice, fmt.Sprintf("Form field '%s' was not a slice", key))
			if okSlice && len(formValueSlice) > 0 {
				assert.Equal(t, expectedValue, formValueSlice[0].(string), fmt.Sprintf("Form field '%s' value mismatch", key))
			} else if okSlice {
				assert.Fail(t, fmt.Sprintf("Form field '%s' was an empty slice", key))
			}
		}
	}

	if expectedRawData != "" {
		assert.Equal(t, expectedRawData, strings.TrimSpace(postResp.Data), "Raw data in response mismatch")
	}


	assert.Nil(t, resp.Bytes(), "JSON decoding should fully consume the response stream (Bytes)")
	}

	if resp.String() != "" {
		assert.Fail(t, "JSON decoding did not fully consume the response stream (String)", resp.String())
	}

	if resp.StatusCode != 200 {
		assert.Fail(t, "Response returned a non-200 code")
	}

	return myJSONStruct
}

func TestPostInvalidURLSession(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Post("%../dir/", nil); err == nil {
		assert.Fail(t, "Some how the request was valid to make request ", err)
	}
}
