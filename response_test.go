package grequests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseOk(t *testing.T) {
	status := []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}
	for _, status := range status {
		verifyResponseOkForStatus(status, t)
	}
}

func verifyResponseOkForStatus(status int, t *testing.T) {
	url := "http://httpbin.org/status/" + strconv.Itoa(status)
	resp, err := Get(url)

	if err != nil {
		assert.Fail(t, "Unable to make request", err)
	}

	if resp.Ok != true {
		assert.Fail(t, fmt.Sprintf("Request did not return OK. Received status code %d rather a 2xx.", resp.StatusCode))
	}
}

func TestReadAfterString(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Hello, world!")
	}))
	defer ts.Close()

	resp, err := Get(ts.URL, nil)
	assert.NoError(t, err, "Expected no error from Get")
	assert.NotNil(t, resp, "Response should not be nil")

	// Call resp.String() to populate internal buffer and close raw response body
	str := resp.String()
	assert.Equal(t, "Hello, world!", str, "String content does not match")

	// Now call resp.Read()
	buf := make([]byte, 10)
	n, err := resp.Read(buf)

	// Assert that it returns an error or 0 bytes read with io.EOF
	// Depending on the exact implementation, reading a closed body might return io.EOF or another error.
	// A successful read (n > 0 and err == nil) would be incorrect.
	if err == nil && n > 0 {
		t.Errorf("Read after String should have failed or returned 0 bytes with io.EOF, but got %d bytes and no error", n)
	} else if err != nil {
		// This is an acceptable outcome, an error occurred as expected.
		// We can optionally check for specific errors, e.g., if it should be exactly io.EOF
		// or a custom error indicating the body is closed. For now, any error is fine.
		t.Logf("Read after String returned expected error: %v", err)
	} else if n == 0 && err == nil { // Some implementations might return 0, nil if already EOF
		t.Logf("Read after String returned 0 bytes and no error, likely EOF.")
	}
	// If n == 0 and err == io.EOF, that's also fine. The check `err != nil` covers this.
}

func TestStringAfterRead(t *testing.T) {
	// Scenario 1: Raw body fully consumed by Read()
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Full consume")
	}))
	defer ts1.Close()

	resp1, err := Get(ts1.URL, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp1)

	buf1 := make([]byte, 50) // Buffer larger than content
	n1, err1 := io.ReadFull(resp1.RawResponse.Body, buf1[:12]) // Read exactly "Full consume"
	assert.NoError(t, err1, "ReadFull should consume the exact number of bytes")
	assert.Equal(t, 12, n1, "Should have read 12 bytes")
	// Intentionally do not call resp1.Close() here to see how String behaves with an open but fully consumed body.
	// However, typical use of ReadFull would be followed by a Close.
	// For robustness, String() should ideally still work or return empty.
	// After ReadFull, the pointer in RawResponse.Body is at EOF.

	str1 := resp1.String()
	// If populateResponseByteBuffer tries to read from RawResponse.Body after it's at EOF,
	// it should get 0 bytes, and the internal buffer should remain empty or reflect prior partial reads.
	// Given String() calls populateResponseByteBuffer which reads from r.RawResponse.Body,
	// and if r.RawResponse.Body is already at EOF, str1 should be empty.
	assert.Equal(t, "", str1, "String after full Read should be empty")
	resp1.RawResponse.Body.Close() // Clean up

	// Scenario 2: Raw body partially consumed by Read()
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Partial consume then String")
	}))
	defer ts2.Close()

	resp2, err := Get(ts2.URL, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)

	buf2 := make([]byte, 8) // Read "Partial "
	n2, err2 := resp2.Read(buf2)
	assert.NoError(t, err2)
	assert.Equal(t, 8, n2)
	assert.Equal(t, "Partial ", string(buf2[:n2]))

	str2 := resp2.String()
	// String() should now read the rest of the body: "consume then String"
	assert.Equal(t, "consume then String", str2, "String after partial Read should return remaining data")
	resp2.RawResponse.Body.Close() // Clean up

	// Scenario 3: Raw body not consumed at all before String()
	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "No read then String")
	}))
	defer ts3.Close()

	resp3, err := Get(ts3.URL, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp3)

	str3 := resp3.String()
	assert.Equal(t, "No read then String", str3, "String with no prior Read should return full data")
	resp3.RawResponse.Body.Close() // Clean up
}

func TestDownloadToFileAfterString(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Download content")
	}))
	defer ts.Close()

	resp, err := Get(ts.URL, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Call String() to populate the internal buffer
	str := resp.String()
	assert.Equal(t, "Download content", str)

	// Download to file
	tempFile, err := os.CreateTemp("", "downloadTest_*.txt")
	assert.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tempFile.Name()) // Clean up

	err = resp.DownloadToFile(tempFile.Name())
	assert.NoError(t, err, "DownloadToFile after String failed")

	// Verify content of downloaded file
	fileContent, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err, "Failed to read downloaded file")
	assert.Equal(t, "Download content", string(fileContent), "Downloaded file content mismatch")

	// Ensure the original response body is closed after DownloadToFile,
	// though String() would have already closed it.
	// Attempting to read again from RawResponse.Body should fail or return EOF.
	// This also tests if DownloadToFile correctly uses the buffer if String was called.
	_, err = resp.RawResponse.Body.Read(make([]byte, 1))
	assert.Error(t, err, "RawResponse.Body should be closed after String() and DownloadToFile")
}

func TestStringAfterDownloadToFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Content for download then string")
	}))
	defer ts.Close()

	resp, err := Get(ts.URL, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	tempFile, err := os.CreateTemp("", "downloadToStringTest_*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Call DownloadToFile first
	err = resp.DownloadToFile(tempFile.Name())
	assert.NoError(t, err, "DownloadToFile failed")

	// Verify file content just to be sure
	fileData, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "Content for download then string", string(fileData))

	// Now call String()
	// Since DownloadToFile consumes the response body and doesn't typically populate the internalByteBuffer,
	// String() should find the raw body closed/EOF and the internal buffer empty.
	str := resp.String()
	assert.Equal(t, "", str, "String after DownloadToFile should be empty")
}

func TestJSONAfterString(t *testing.T) {
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	jsonData := `{"name":"TestJSONAfterString","value":123}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, jsonData)
	}))
	defer ts.Close()

	resp, err := Get(ts.URL, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Call String() to populate the internal buffer
	strContent := resp.String()
	assert.Equal(t, jsonData, strContent)

	// Now call JSON()
	var data TestData
	err = resp.JSON(&data)
	assert.NoError(t, err, "JSON() after String failed")

	// Verify parsed data
	assert.Equal(t, "TestJSONAfterString", data.Name)
	assert.Equal(t, 123, data.Value)

	// The raw response body should be closed by String()
	// JSON() should have used the internal buffer.
	_, err = resp.RawResponse.Body.Read(make([]byte, 1))
	assert.Error(t, err, "RawResponse.Body should be closed after String() and JSON()")
}

func TestStringAfterJSON(t *testing.T) {
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	jsonData := `{"name":"TestStringAfterJSON","value":456}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, jsonData)
	}))
	defer ts.Close()

	resp, err := Get(ts.URL, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Call JSON() first
	var data TestData
	err = resp.JSON(&data)
	assert.NoError(t, err, "JSON() call failed")

	// Verify parsed data
	assert.Equal(t, "TestStringAfterJSON", data.Name)
	assert.Equal(t, 456, data.Value)

	// Now call String()
	// JSON() consumes the response body and typically doesn't populate the internalByteBuffer for String()
	// unless explicitly designed to do so (which is not the standard behavior being tested).
	// Thus, String() should find the raw body closed/EOF and the internal buffer empty.
	strContent := resp.String()
	assert.Equal(t, "", strContent, "String after JSON() should be empty")
}
