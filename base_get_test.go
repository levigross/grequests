package grequests

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

type BasicGetResponse struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept          string `json:"Accept"`
		Accept_Encoding string `json:"Accept-Encoding"`
		Accept_Language string `json:"Accept-Language"`
		Dnt             string `json:"Dst"`
		Host            string `json:"Host"`
		User_Agent      string `json:"User-Agent"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseNewHeader struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept             string `json:"Accept"`
		Accept_Encoding    string `json:"Accept-Encoding"`
		Accept_Language    string `json:"Accept-Language"`
		Dnt                string `json:"Dst"`
		Host               string `json:"Host"`
		User_Agent         string `json:"User-Agent"`
		X_Wonderful_Header string `json:"X-Wonderful-Header"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseBasicAuth struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept          string `json:"Accept"`
		Accept_Encoding string `json:"Accept-Encoding"`
		Accept_Language string `json:"Accept-Language"`
		Dnt             string `json:"Dst"`
		Host            string `json:"Host"`
		User_Agent      string `json:"User-Agent"`
		Authorization   string `json:"Authorization"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseArgs struct {
	Args struct {
		Goodbye string `json:"Goodbye"`
		Hello   string `json:"Hello"`
	} `json:"args"`
	Headers struct {
		Accept          string `json:"Accept"`
		Accept_Encoding string `json:"Accept-Encoding"`
		Accept_Language string `json:"Accept-Language"`
		Dnt             string `json:"Dst"`
		Host            string `json:"Host"`
		User_Agent      string `json:"User-Agent"`
		Authorization   string `json:"Authorization"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

func TestGetNoOptions(t *testing.T) {
	verifyOkResponse(<-Get("http://httpbin.org/get", nil), t)
}

func TestGetNoOptionsChannel(t *testing.T) {
	respChan := Get("http://httpbin.org/get", nil)
	select {
	case resp := <-respChan:
		verifyOkResponse(resp, t)
	}
}

func TestGetCustomUserAgent(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1"}
	resp := <-Get("http://httpbin.org/get", ro)
	jsonResp := verifyOkResponse(resp, t)
	if jsonResp.Headers.User_Agent != "LeviBot 0.1" {
		t.Error("User agent header not properly set")
	}
}

func TestGetBasicAuth(t *testing.T) {
	ro := &RequestOptions{Auth: []string{"Levi", "Bot"}}
	resp := <-Get("http://httpbin.org/get", ro)
	// Not the usual JSON so copy and paste from below

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJsonStruct := &BasicGetResponseBasicAuth{}

	err := resp.Json(myJsonStruct)
	if err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJsonStruct.Headers.Authorization != "Basic TGV2aTpCb3Q=" {
		t.Error("Unable to set HTTP basic auth", myJsonStruct.Headers)
	}

}

func TestGetCustomHeader(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1",
		Headers: map[string]string{"X-Wonderful-Header": "1"}}
	resp := <-Get("http://httpbin.org/get", ro)
	// Not the usual JSON so copy and paste from below

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJsonStruct := &BasicGetResponseNewHeader{}

	err := resp.Json(myJsonStruct)
	if err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJsonStruct.Headers.X_Wonderful_Header != "1" {
		t.Error("Unable to set custom HTTP header", myJsonStruct.Headers)
	}
}

func TestGetInvalidSSLCertNoVerify(t *testing.T) {
	ro := &RequestOptions{InsecureSkipVerify: true}
	resp := <-Get("https://www.pcwebshop.co.uk/", ro)
	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}
}

func TestGetInvalidSSLCert(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1"}
	resp := <-Get("https://www.pcwebshop.co.uk/", ro)

	if resp.Error == nil {
		t.Error("SSL verification worked when it shouldn't of", resp.Error)
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

}

func TestGetBasicArgs(t *testing.T) {
	ro := &RequestOptions{
		Params: map[string]string{"Hello": "World"},
	}
	verifyOkArgsResponse(<-Get("http://httpbin.org/get?Goodbye=World", ro), t)

}

func TestGetBasicArgsParams(t *testing.T) {
	ro := &RequestOptions{
		Params: map[string]string{"Hello": "World", "Goodbye": "World"},
	}
	verifyOkArgsResponse(<-Get("http://httpbin.org/get", ro), t)
}

func TestGetBasicArgsParamsOverwrite(t *testing.T) {
	ro := &RequestOptions{
		Params: map[string]string{"Hello": "World", "Goodbye": "World"},
	}
	verifyOkArgsResponse(<-Get("http://httpbin.org/get?Hello=Nothing", ro), t)
}

func TestGetFileDownload(t *testing.T) {
	resp := <-Get("http://httpbin.org/get", nil)

	fileName := "randomFile"

	if err := resp.DownloadToFile(fileName); err != nil {
		t.Error("Unable to download to file: ", err)
	}

	fd, err := os.Open(fileName)
	defer fd.Close()
	defer os.Remove(fileName)

	if err != nil {
		t.Error("Unable to open file to verify content ", err)
	}

	jsonDecoder := json.NewDecoder(fd)

	myJsonStruct := &BasicGetResponse{}

	if err := jsonDecoder.Decode(myJsonStruct); err != nil {
		t.Error("Unable to cocerce file to JSON ", err)
	}

	if myJsonStruct.URL != "http://httpbin.org/get" {
		t.Error("For some reason the URL isn't the same", myJsonStruct.URL)
	}

	if myJsonStruct.Headers.Host != "httpbin.org" {
		t.Error("The host header is invalid")
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

}

func TestGetBytes(t *testing.T) {
	resp := <-Get("http://httpbin.org/get", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.Bytes() == nil {
		t.Error("JSON decoding did not fully consume the response stream")
	}

	if bytes.Compare(resp.Bytes(), resp.Bytes()) != 0 {
		t.Error("Body bytes have not been cached", resp.Bytes())
	}
}

func TestGetBytesNoBuffer(t *testing.T) {
	resp := <-Get("http://httpbin.org/get", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.Bytes() == nil {
		t.Error("JSON decoding did not fully consume the response stream")
	}

	if bytes.Compare(resp.Bytes(), resp.Bytes()) != 0 {
		t.Error("Body bytes have not been cached", resp.Bytes())
	}

	resp.ClearInternalBuffer()

	if resp.Bytes() != nil {
		t.Error("Internal Buffer not cleaned up")
	}
}

func TestGetString(t *testing.T) {
	resp := <-Get("http://httpbin.org/get", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.String() == "" {
		t.Error("JSON decoding did not fully consume the response stream (string)", resp.String())
	}

	if resp.String() != resp.String() {
		t.Error("Body string have not been cached", resp.String())
	}

	if err := resp.DownloadToFile("randomFile"); err != nil {
		t.Error("Unable to download file: ", err)
	}

	defer os.Remove("randomFile")

}

func verifyOkArgsResponse(resp *Response, t *testing.T) *BasicGetResponseArgs {
	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJsonStruct := &BasicGetResponseArgs{}

	err := resp.Json(myJsonStruct)
	if err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJsonStruct.Args.Goodbye != "World" && myJsonStruct.Args.Hello != "World" {
		t.Error("Args not properly set", myJsonStruct.Args)
	}

	if myJsonStruct.URL != "http://httpbin.org/get?Goodbye=World&Hello=World" {
		t.Error("Url is not properly constructed", myJsonStruct.URL)
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

	return myJsonStruct
}

// verifyResponse will verify the following conditions
// 1. The request didn't return an error
// 2. The response returned an OK (a status code within the 200 range)
// 3. The output can be coerced to JSON (this may change later)
// It should only be run when testing GET request to http://httpbin.org/get expecting JSON
func verifyOkResponse(resp *Response, t *testing.T) *BasicGetResponse {
	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJsonStruct := &BasicGetResponse{}

	err := resp.Json(myJsonStruct)
	if err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJsonStruct.URL != "http://httpbin.org/get" {
		t.Error("For some reason the URL isn't the same", myJsonStruct.URL)
	}

	if myJsonStruct.Headers.Host != "httpbin.org" {
		t.Error("The host header is invalid")
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

	return myJsonStruct
}
