package grequests_test

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/levigross/grequests"
)

func Example_basic_get() {
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

func Example_Cookies() {
	resp, err := grequests.Get("http://httpbin.org/cookies",
		&grequests.RequestOptions{
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

	if err != nil {
		log.Println("Unable to make request", err)
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
	}

	log.Println(resp.String())
}

func Example_Session() {
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

func Example_Parse_XML() {
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

func xmlASCIIDecoder(charset string, input io.Reader) (io.Reader, error) {
	return input, nil
}
