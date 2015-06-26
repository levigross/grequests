package grequests

import (
	"encoding/xml"
	"math"
	"strings"
	"testing"
)

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

func TestBasicPostRequest(t *testing.T) {
	verifyOkPostResponse(<-PostAsync("http://httpbin.org/post",
		&RequestOptions{Data: map[string]string{"One": "Two"}}), t)

}

func TestBasicRegularPostRequest(t *testing.T) {
	verifyOkPostResponse(Post("http://httpbin.org/post",
		&RequestOptions{Data: map[string]string{"One": "Two"}}), t)

}

func TestBasicPostRequestInvalidURL(t *testing.T) {
	resp := <-PostAsync("%../dir/",
		&RequestOptions{Data: map[string]string{"One": "Two"},
			Params: map[string]string{"1": "2"}})

	if resp.Error == nil {
		t.Error("Somehow the request went through")
	}

}

func TestBasicPostRequestInvalidURLNoParams(t *testing.T) {
	resp := <-PostAsync("%../dir/", &RequestOptions{Data: map[string]string{"One": "Two"}})

	if resp.Error == nil {
		t.Error("Somehow the request went through")
	}

}

func TestXMLPostRequestInvalidURL(t *testing.T) {
	resp := <-PostAsync("%../dir/",
		&RequestOptions{XML: XMLPostMessage{Name: "Human", Age: 1, Height: 1}})

	if resp.Error == nil {
		t.Error("Somehow the request went through")
	}
}

func TestBasicPostJsonRequestInvalidURL(t *testing.T) {
	resp := Post("%../dir/",
		&RequestOptions{JSON: map[string]string{"One": "Two"}, IsAjax: true})

	if resp.Error == nil {
		t.Error("Somehow the request went through")
	}
}

func TestBasicPostJsonRequestInvalidJSON(t *testing.T) {
	resp := Post("http://httpbin.org/post",
		&RequestOptions{JSON: math.NaN(), IsAjax: true})

	if resp.Error == nil {
		t.Error("Somehow the request went through")
	}

	if resp.Ok == true {
		t.Error("Somehow the request is OK")
	}
}

func TestBasicPostJsonRequestInvalidXML(t *testing.T) {
	resp := Post("http://httpbin.org/post",
		&RequestOptions{XML: map[string]string{"One": "two"}, IsAjax: true})

	if resp.Error == nil {
		t.Error("Somehow the request went through")
	}

	if resp.Ok == true {
		t.Error("Somehow the request is OK")
	}
}

func TestBasicPostRequestUploadInvalidURL(t *testing.T) {

	fd, err := FileUploadFromDisk("test_files/mypassword")

	if err != nil {
		t.Error("Unable to open file: ", err)
	}

	defer fd.FileContents.Close()

	resp := <-PostAsync("%../dir/",
		&RequestOptions{
			File: fd,
			Data: map[string]string{"One": "Two"},
		})

	if resp.Error == nil {
		t.Fatal("Somehow able to make the request")
	}
}

func TestBasicPostRequestUploadInvalidFileUpload(t *testing.T) {

	resp := <-PostAsync("%../dir/",
		&RequestOptions{
			File: &FileUpload{FileName: "üfdsufhidÄDJSHAKÔÓÔ", FileContents: nil},
			Data: map[string]string{"One": "Two"},
		})

	if resp.Error == nil {
		t.Fatal("Somehow able to make the request")
	}
}

func TestXMLPostRequest(t *testing.T) {
	resp := <-PostAsync("http://httpbin.org/post",
		&RequestOptions{XML: XMLPostMessage{Name: "Human", Age: 1, Height: 1}})

	if resp.Error != nil {
		t.Fatal("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicPostJSONResponse{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	myXMLStruct := &XMLPostMessage{}

	xml.Unmarshal([]byte(myJSONStruct.Data), myXMLStruct)

	if myXMLStruct.Age != 1 {
		t.Errorf("Unable to serialize XML response from within JSON %#v ", myXMLStruct)
	}

}

func TestBasicPostRequestUpload(t *testing.T) {

	fd, err := FileUploadFromDisk("test_files/mypassword")

	if err != nil {
		t.Error("Unable to open file: ", err)
	}

	resp := <-PostAsync("http://httpbin.org/post",
		&RequestOptions{
			File: fd,
			Data: map[string]string{"One": "Two"},
		})

	if resp.Error != nil {
		t.Fatal("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicPostFileUpload{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.URL != "http://httpbin.org/post" {
		t.Error("For some reason the URL isn't the same", myJSONStruct.URL)
	}

	if myJSONStruct.Headers.Host != "httpbin.org" {
		t.Error("The host header is invalid")
	}

	if myJSONStruct.Files.File != "saucy sauce" {
		t.Error("File upload contents have been modified: ", myJSONStruct.Files.File)
	}

	if resp.Bytes() != nil {
		t.Error("JSON decoding did not fully consume the response stream (Bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("JSON decoding did not fully consume the response stream (String)", resp.String())
	}

	if resp.StatusCode != 200 {
		t.Error("Response returned a non-200 code")
	}

	if myJSONStruct.Form.One != "Two" {
		t.Error("Unable to parse form properly", myJSONStruct.Form)
	}

}

func TestBasicPostJsonRequest(t *testing.T) {
	resp := <-PostAsync("http://httpbin.org/post",
		&RequestOptions{JSON: map[string]string{"One": "Two"}, IsAjax: true})

	if resp.Error != nil {
		t.Fatal("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicPostJSONResponse{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.URL != "http://httpbin.org/post" {
		t.Error("For some reason the URL isn't the same", myJSONStruct.URL)
	}

	if myJSONStruct.Headers.Host != "httpbin.org" {
		t.Error("The host header is invalid")
	}

	if myJSONStruct.JSON.One != "Two" {
		t.Error("Invalid post response: ", myJSONStruct.JSON.One)
	}

	if strings.TrimSpace(myJSONStruct.Data) != `{"One":"Two"}` {
		t.Error("JSON not properly returned: ", myJSONStruct.Data)
	}

	if resp.Bytes() != nil {
		t.Error("JSON decoding did not fully consume the response stream (Bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("JSON decoding did not fully consume the response stream (String)", resp.String())
	}

	if resp.StatusCode != 200 {
		t.Error("Response returned a non-200 code")
	}

	if myJSONStruct.Headers.XRequestedWith != "XMLHttpRequest" {
		t.Error("Invalid requested header: ", myJSONStruct.Headers.XRequestedWith)
	}

}

// verifyResponse will verify the following conditions
// 1. The request didn't return an error
// 2. The response returned an OK (a status code within the 200 range)
// 3. The output can be coerced to JSON (this may change later)
// It should only be run when testing GET request to http://httpbin.org/post expecting JSON
func verifyOkPostResponse(resp *Response, t *testing.T) *BasicPostResponse {
	if resp.Error != nil {
		t.Fatal("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicPostResponse{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.URL != "http://httpbin.org/post" {
		t.Error("For some reason the URL isn't the same", myJSONStruct.URL)
	}

	if myJSONStruct.Headers.Host != "httpbin.org" {
		t.Error("The host header is invalid")
	}

	if myJSONStruct.Form.One != "Two" {
		t.Errorf("Invalid post response: %#v", myJSONStruct.Form)
	}

	if resp.Bytes() != nil {
		t.Error("JSON decoding did not fully consume the response stream (Bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("JSON decoding did not fully consume the response stream (String)", resp.String())
	}

	if resp.StatusCode != 200 {
		t.Error("Response returned a non-200 code")
	}

	return myJSONStruct
}
