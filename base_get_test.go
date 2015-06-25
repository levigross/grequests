package grequests

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"testing"
)

type BasicGetResponse struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
		Dnt            string `json:"Dst"`
		Host           string `json:"Host"`
		UserAgent      string `json:"User-Agent"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseNewHeader struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept           string `json:"Accept"`
		AcceptEncoding   string `json:"Accept-Encoding"`
		AcceptLanguage   string `json:"Accept-Language"`
		Dnt              string `json:"Dst"`
		Host             string `json:"Host"`
		UserAgent        string `json:"User-Agent"`
		XWonderfulHeader string `json:"X-Wonderful-Header"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseBasicAuth struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
		Dnt            string `json:"Dst"`
		Host           string `json:"Host"`
		UserAgent      string `json:"User-Agent"`
		Authorization  string `json:"Authorization"`
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
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
		Dnt            string `json:"Dst"`
		Host           string `json:"Host"`
		UserAgent      string `json:"User-Agent"`
		Authorization  string `json:"Authorization"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type GetXMLSample struct {
	XMLName xml.Name `xml:"slideshow"`
	Title   string   `xml:"title,attr"`
	Date    string   `xml:"date,attr"`
	Author  string   `xml:"author,attr"`
	Slide   []struct {
		Type  string `xml:"type,attr"`
		Title string `xml:"title"`
	} `xml:"slide"`
}

type TestJSONCookies struct {
	Cookies struct {
		AnotherCookie string `json:"AnotherCookie"`
		TestCookie    string `json:"TestCookie"`
	} `json:"cookies"`
}

func TestGetNoOptions(t *testing.T) {
	verifyOkResponse(<-GetAsync("http://httpbin.org/get", nil), t)
}

func TestGetSyncNoOptions(t *testing.T) {
	verifyOkResponse(Get("http://httpbin.org/get", nil), t)
}

func TestGetNoOptionsGzip(t *testing.T) {
	verifyOkResponse(<-GetAsync("https://httpbin.org/gzip", nil), t)
}

func TestGetWithCookies(t *testing.T) {
	resp := <-GetAsync("http://httpbin.org/cookies",
		&RequestOptions{
			Cookies: []http.Cookie{
				http.Cookie{
					Name:     "TestCookie",
					Value:    "Random Value",
					HttpOnly: true,
					Secure:   false,
				}, http.Cookie{
					Name:     "AnotherCookie",
					Value:    "Some Value",
					HttpOnly: true,
					Secure:   false,
				},
			},
		})

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &TestJSONCookies{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Cannot serialize cookie JSON: ", err)
	}

	if myJSONStruct.Cookies.TestCookie != "Random Value" {
		t.Errorf("Cookie value not set properly: %#v", myJSONStruct)
	}

	if myJSONStruct.Cookies.AnotherCookie != "Some Value" {
		t.Errorf("Cookie value not set properly: %#v", myJSONStruct)
	}

}

//func TestGetNoOptionsDeflate(t *testing.T) {
//	verifyOkResponse(<-GetAsync("http://httpbin.org/deflate", nil), t)
//}

func xmlASCIIDecoder(charset string, input io.Reader) (io.Reader, error) {
	return input, nil
}

func TestGetInvalidURL(t *testing.T) {
	resp := <-GetAsync("%../dir/", &RequestOptions{Params: map[string]string{"1": "2"}})

	if resp.Error == nil {
		t.Error("Some how the request was valid to make request", resp.Error)
	}
}

func TestGetInvalidURLNoParams(t *testing.T) {
	resp := <-GetAsync("%../dir/", nil)

	if resp.Error == nil {
		t.Error("Some how the request was valid to make request", resp.Error)
	}
}

func TestGetXMLSerialize(t *testing.T) {
	resp := <-GetAsync("http://httpbin.org/xml", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	userXML := &GetXMLSample{}

	if err := resp.XML(userXML, xmlASCIIDecoder); err != nil {
		t.Error("Unable to consume the response as XML: ", err)
	}

	if userXML.Title != "Sample Slide Show" {
		t.Errorf("Invalid XML serialization %#v", userXML)
	}

	if err := resp.XML(int(123), nil); err == nil {
		t.Error("Still able to consume XML from used response")
	}

}

func TestGetNoOptionsChannel(t *testing.T) {
	respChan := GetAsync("http://httpbin.org/get", nil)
	select {
	case resp := <-respChan:
		verifyOkResponse(resp, t)
	}
}

func TestGetCustomUserAgent(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1"}
	resp := <-GetAsync("http://httpbin.org/get", ro)
	jsonResp := verifyOkResponse(resp, t)
	if jsonResp.Headers.UserAgent != "LeviBot 0.1" {
		t.Error("User agent header not properly set")
	}
}

func TestGetBasicAuth(t *testing.T) {
	ro := &RequestOptions{Auth: []string{"Levi", "Bot"}}
	resp := <-GetAsync("http://httpbin.org/get", ro)
	// Not the usual JSON so copy and paste from below

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicGetResponseBasicAuth{}

	err := resp.JSON(myJSONStruct)
	if err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.Headers.Authorization != "Basic TGV2aTpCb3Q=" {
		t.Error("Unable to set HTTP basic auth", myJSONStruct.Headers)
	}

}

func TestGetCustomHeader(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1",
		Headers: map[string]string{"X-Wonderful-Header": "1"}}
	resp := <-GetAsync("http://httpbin.org/get", ro)
	// Not the usual JSON so copy and paste from below

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicGetResponseNewHeader{}

	err := resp.JSON(myJSONStruct)
	if err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.Headers.XWonderfulHeader != "1" {
		t.Error("Unable to set custom HTTP header", myJSONStruct.Headers)
	}
}

func TestGetInvalidSSLCertNoVerify(t *testing.T) {
	ro := &RequestOptions{InsecureSkipVerify: true}
	resp := <-GetAsync("https://www.pcwebshop.co.uk/", ro)
	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}
}

func TestGetInvalidSSLCertNoVerifyNoOptions(t *testing.T) {
	resp := <-GetAsync("https://www.pcwebshop.co.uk/", nil)
	if resp.Error == nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok == true {
		t.Error("Request did not return OK")
	}
}

func TestGetInvalidSSLCertNoCompression(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1", DisableCompression: true}
	resp := <-GetAsync("https://www.pcwebshop.co.uk/", ro)

	if resp.Error == nil {
		t.Error("SSL verification worked when it shouldn't of", resp.Error)
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

}

func TestGetInvalidSSLCertWithCompression(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1", DisableCompression: false}
	resp := <-GetAsync("https://www.pcwebshop.co.uk/", ro)

	if resp.Error == nil {
		t.Error("SSL verification worked when it shouldn't of", resp.Error)
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

}

func TestErrorResponseNOOP(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1", DisableCompression: false}
	resp := <-GetAsync("https://www.pcwebshop.co.uk/", ro)

	if resp.Error == nil {
		t.Error("SSL verification worked when it shouldn't of", resp.Error)
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

	myJSONStruct := &BasicGetResponseArgs{}

	if err := resp.JSON(myJSONStruct); err == nil {
		t.Error("Somehow Able to convert to JSON", err)
	}

	if resp.Bytes() != nil {
		t.Error("Somehow byte buffer is working now (bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("Somehow byte buffer is working now (bytes)", resp.String())
	}

	resp.ClearInternalBuffer()

	if resp.Bytes() != nil {
		t.Error("Somehow byte buffer is working now (bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("Somehow byte buffer is working now (bytes)", resp.String())
	}

	userXML := &GetXMLSample{}

	if err := resp.XML(userXML, xmlASCIIDecoder); err == nil {
		t.Errorf("Somehow to consume the response as XML: %#v", userXML)
	}

	fileName := "randomFile"

	if err := resp.DownloadToFile(fileName); err == nil {
		t.Error("Somehow able to download to file: ", err)
	}

	var buf [1]byte

	if written, err := resp.Read(buf[:]); written != -1 && err == nil {
		t.Error("Somehow we were able to read from our error response")
	}

}

func TestGetInvalidSSLCertNoCompressionNoVerify(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1", InsecureSkipVerify: true, DisableCompression: true}
	resp := <-GetAsync("https://www.pcwebshop.co.uk/", ro)

	if resp.Error != nil {
		t.Error("SSL verification worked when it shouldn't of", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did return OK")
	}

}

func TestGetInvalidSSLCertWithCompressionNoVerify(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1", InsecureSkipVerify: true, DisableCompression: false}
	resp := <-GetAsync("https://www.pcwebshop.co.uk/", ro)

	if resp.Error != nil {
		t.Error("SSL verification worked when it shouldn't of", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did return OK")
	}

}

func TestGetInvalidSSLCert(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1"}
	resp := <-GetAsync("https://www.pcwebshop.co.uk/", ro)

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
	verifyOkArgsResponse(<-GetAsync("http://httpbin.org/get?Goodbye=World", ro), t)

}

func TestGetBasicArgsParams(t *testing.T) {
	ro := &RequestOptions{
		Params: map[string]string{"Hello": "World", "Goodbye": "World"},
	}
	verifyOkArgsResponse(<-GetAsync("http://httpbin.org/get", ro), t)
}

func TestGetBasicArgsParamsOverwrite(t *testing.T) {
	ro := &RequestOptions{
		Params: map[string]string{"Hello": "World", "Goodbye": "World"},
	}
	verifyOkArgsResponse(<-GetAsync("http://httpbin.org/get?Hello=Nothing", ro), t)
}

func TestGetFileDownload(t *testing.T) {
	resp := <-GetAsync("http://httpbin.org/get", nil)

	fileName := "randomFile"

	if err := resp.DownloadToFile(fileName); err != nil {
		t.Error("Unable to download to file: ", err)
	}

	if err := resp.DownloadToFile("."); err == nil {
		t.Error("Able to create file '.'")
	}

	fd, err := os.Open(fileName)
	defer fd.Close()
	defer os.Remove(fileName)

	if err != nil {
		t.Error("Unable to open file to verify content ", err)
	}

	jsonDecoder := json.NewDecoder(fd)

	myJSONStruct := &BasicGetResponse{}

	if err := jsonDecoder.Decode(myJSONStruct); err != nil {
		t.Error("Unable to cocerce file to JSON ", err)
	}

	if myJSONStruct.URL != "http://httpbin.org/get" {
		t.Error("For some reason the URL isn't the same", myJSONStruct.URL)
	}

	if myJSONStruct.Headers.Host != "httpbin.org" {
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

func TestJsonConsumedResponse(t *testing.T) {
	resp := <-GetAsync("http://httpbin.org/get", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.Bytes() == nil {
		t.Error("Unable to coerce value to bytes", resp.Bytes())
	}

	resp.ClearInternalBuffer()

	if err := resp.JSON(struct{}{}); err == nil {
		t.Error("Struct should not be able to hold JSON: ")
	}
}

func TestDownloadConsumedResponse(t *testing.T) {
	resp := <-GetAsync("http://httpbin.org/get", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.Bytes() == nil {
		t.Error("Unable to coerce value to bytes")
	}

	resp.ClearInternalBuffer()

	if err := resp.DownloadToFile("randomFile"); err == nil {
		t.Error("Still able to download file: ", err)
	}

	defer os.Remove("randomFile")
}

func TestGetBytes(t *testing.T) {
	resp := <-GetAsync("http://httpbin.org/get", nil)

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
	resp := <-GetAsync("http://httpbin.org/get", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.Bytes() == nil {
		t.Error("Cannot coerce HTTP response to bytes")
	}

	if bytes.Compare(resp.Bytes(), resp.Bytes()) != 0 {
		t.Error("Body bytes have not been cached", resp.Bytes())
	}

	if err := resp.DownloadToFile("randomFile"); err != nil {
		t.Error("Unable to download file: ", err)
	}

	defer os.Remove("randomFile")

	resp.ClearInternalBuffer()

	if resp.Bytes() != nil {
		t.Error("Internal Buffer not cleaned up")
	}
}

func TestGetString(t *testing.T) {
	resp := <-GetAsync("http://httpbin.org/get", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.String() == "" {
		t.Error("Response Stream not returned as string", resp.String())
	}

	if resp.String() != resp.String() {
		t.Error("Body string have not been cached", resp.String())
	}

	if err := resp.DownloadToFile("randomFile"); err != nil {
		t.Error("Unable to download file: ", err)
	}

	defer os.Remove("randomFile")

	resp.ClearInternalBuffer()

	if resp.String() != "" {
		t.Error("Internal Buffer not cleaned up")
	}

}

func verifyOkArgsResponse(resp *Response, t *testing.T) *BasicGetResponseArgs {
	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicGetResponseArgs{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.Args.Goodbye != "World" && myJSONStruct.Args.Hello != "World" {
		t.Error("Args not properly set", myJSONStruct.Args)
	}

	if myJSONStruct.URL != "http://httpbin.org/get?Goodbye=World&Hello=World" {
		t.Error("Url is not properly constructed", myJSONStruct.URL)
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

	myJSONStruct := &BasicGetResponse{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.Headers.Host != "httpbin.org" {
		t.Error("The host header is invalid")
	}

	if resp.Bytes() != nil {
		t.Errorf("JSON decoding did not fully consume the response stream (Bytes) %#v", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("JSON decoding did not fully consume the response stream (String)", resp.String())
	}

	if resp.StatusCode != 200 {
		t.Error("Response returned a non-200 code")
	}

	return myJSONStruct
}
