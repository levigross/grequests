package grequests_test

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/levigross/grequests"
)

func Example_basicGet() {
	// This is a very basic GET request
	resp, err := grequests.Get("http://httpbin.org/get", nil)

	if err != nil {
		log.Println(err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}

	log.Println(resp.String())
}

func Example_basicGetCustomHTTPClient() {
	// This is a very basic GET request
	resp, err := grequests.Get("http://httpbin.org/get", &grequests.RequestOptions{HTTPClient: http.DefaultClient})

	if err != nil {
		log.Println(err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}

	log.Println(resp.String())
}

func Example_proxy() {
	proxyURL, err := url.Parse("http://127.0.0.1:8080") // Proxy URL
	if err != nil {
		log.Panicln(err)
	}

	resp, err := grequests.Get("http://www.levigross.com/",
		&grequests.RequestOptions{Proxies: map[string]*url.URL{proxyURL.Scheme: proxyURL}})

	if err != nil {
		log.Println(err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}

	log.Println(resp)
}

func Example_cookies() {
	resp, err := grequests.Get("http://httpbin.org/cookies",
		&grequests.RequestOptions{
			Cookies: []*http.Cookie{
				{
					Name:     "TestCookie",
					Value:    "Random Value",
					HttpOnly: true,
					Secure:   false,
				}, {
					Name:     "AnotherCookie",
					Value:    "Some Value",
					HttpOnly: true,
					Secure:   false,
				},
			},
		})

	if err != nil {
		log.Println("Unable to make request", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}

	log.Println(resp.String())
}

func Example_session() {
	session := grequests.NewSession(nil)

	resp, err := session.Get("http://httpbin.org/cookies/set", &grequests.RequestOptions{Params: map[string]string{"one": "two"}})

	if err != nil {
		log.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}

	log.Println(resp.String())

}

func Example_parse_XML() {
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

	resp, err := grequests.Get("http://httpbin.org/xml", nil)

	if err != nil {
		log.Println("Unable to make request", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}

	userXML := &GetXMLSample{}

	// func xmlASCIIDecoder(charset string, input io.Reader) (io.Reader, error) {
	// 	return input, nil
	// }

	// If the server returns XML encoded in another charset (not UTF-8) – you
	// must provide an encoder function that looks like the one I wrote above.

	// If you an consuming UTF-8 just pass `nil` into the second arg
	if err := resp.XML(userXML, xmlASCIIDecoder); err != nil {
		log.Println("Unable to consume the response as XML: ", err)
	}

	if userXML.Title != "Sample Slide Show" {
		log.Printf("Invalid XML serialization %#v", userXML)
	}
}

func Example_customUserAgent() {
	ro := &grequests.RequestOptions{UserAgent: "LeviBot 0.1"}
	resp, err := grequests.Get("http://httpbin.org/get", ro)

	if err != nil {
		log.Fatal("Oops something went wrong: ", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}

	log.Println(resp.String())
}

func Example_basicAuth() {
	ro := &grequests.RequestOptions{Auth: []string{"Levi", "Bot"}}
	resp, err := grequests.Get("http://httpbin.org/get", ro)
	// Not the usual JSON so copy and paste from below

	if err != nil {
		log.Println("Unable to make request", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}
}

func Example_customHTTPHeader() {
	ro := &grequests.RequestOptions{UserAgent: "LeviBot 0.1",
		Headers: map[string]string{"X-Wonderful-Header": "1"}}
	resp, err := grequests.Get("http://httpbin.org/get", ro)
	// Not the usual JSON so copy and paste from below

	if err != nil {
		log.Println("Unable to make request", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}
}

func Example_acceptInvalidTLSCert() {
	ro := &grequests.RequestOptions{InsecureSkipVerify: true}
	resp, err := grequests.Get("https://www.pcwebshop.co.uk/", ro)

	if err != nil {
		log.Println("Unable to make request", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}
}

func Example_urlQueryParams() {
	ro := &grequests.RequestOptions{
		Params: map[string]string{"Hello": "World", "Goodbye": "World"},
	}
	resp, err := grequests.Get("http://httpbin.org/get", ro)
	// url will now be http://httpbin.org/get?hello=world&goodbye=world

	if err != nil {
		log.Println("Unable to make request", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}
}

func Example_downloadFile() {
	resp, err := grequests.Get("http://httpbin.org/get", nil)

	if err != nil {
		log.Println("Unable to make request", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}

	if err := resp.DownloadToFile("randomFile"); err != nil {
		log.Println("Unable to download to file: ", err)
	}

	if err != nil {
		log.Println("Unable to download file", err)
	}

}

func Example_postForm() {
	resp, err := grequests.Post("http://httpbin.org/post",
		&grequests.RequestOptions{Data: map[string]string{"One": "Two"}})

	// This is the basic form POST. The request body will be `one=two`

	if err != nil {
		log.Println("Cannot post: ", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}
}

func Example_postXML() {

	type XMLPostMessage struct {
		Name   string
		Age    int
		Height int
	}

	resp, err := grequests.Post("http://httpbin.org/post",
		&grequests.RequestOptions{XML: XMLPostMessage{Name: "Human", Age: 1, Height: 1}})
	// The request body will contain the XML generated by the `XMLPostMessage` struct

	if err != nil {
		log.Println("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}
}

func Example_postFileUpload() {

	fd, err := grequests.FileUploadFromDisk("test_files/mypassword")

	if err != nil {
		log.Println("Unable to open file: ", err)
	}

	// This will upload the file as a multipart mime request
	resp, err := grequests.Post("http://httpbin.org/post",
		&grequests.RequestOptions{
			Files: fd,
			Data:  map[string]string{"One": "Two"},
		})

	if err != nil {
		log.Println("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}
}

func Example_postJSONAJAX() {
	resp, err := grequests.Post("http://httpbin.org/post",
		&grequests.RequestOptions{
			JSON:   map[string]string{"One": "Two"},
			IsAjax: true, // this adds the X-Requested-With: XMLHttpRequest header
		})

	if err != nil {
		log.Println("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}

}

func xmlASCIIDecoder(charset string, input io.Reader) (io.Reader, error) {
	return input, nil
}
