package grequests

import (
	"net/url"
	"testing"
)

func TestBasicPutRequest(t *testing.T) {
	resp, err := Put("http://httpbin.org/put", nil)

	if err != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

}

func TestBasicPutUploadRequest(t *testing.T) {
	fd, err := FileUploadFromDisk("testdata/mypassword")

	if err != nil {
		t.Error("Unable to open file: ", err)
	}

	resp, _ := Put("http://httpbin.org/put",
		&RequestOptions{
			Files: fd,
			Data:  map[string]string{"One": "Two"},
		})

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

}

func TestBasicPutUploadRequestInvalidURL(t *testing.T) {
	fd, err := FileUploadFromDisk("testdata/mypassword")

	if err != nil {
		t.Error("Unable to open file: ", err)
	}

	_, err = Put("%../dir/",
		&RequestOptions{
			Files: fd,
			Data:  map[string]string{"One": "Two"},
		})

	if err == nil {
		t.Fatal("Somehow able to make the request")
	}
}

func TestSessionPutUploadRequestInvalidURL(t *testing.T) {
	fd, err := FileUploadFromDisk("testdata/mypassword")

	if err != nil {
		t.Error("Unable to open file: ", err)
	}

	session := NewSession(nil)

	_, err = session.Put("%../dir/",
		&RequestOptions{
			Files: fd,
			Data:  map[string]string{"One": "Two"},
		})

	if err == nil {
		t.Fatal("Somehow able to make the request")
	}
}

func TestPutSession(t *testing.T) {
	session := NewSession(nil)

	resp, err := session.Get("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"one": "two"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	resp, err = session.Get("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"two": "three"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	resp, err = session.Get("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"three": "four"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	resp, err = session.Put("http://httpbin.org/put", &RequestOptions{Data: map[string]string{"one": "two"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	cookieURL, err := url.Parse("http://httpbin.org")
	if err != nil {
		t.Error("We (for some reason) cannot parse the cookie URL")
	}

	if len(session.HTTPClient.Jar.Cookies(cookieURL)) != 3 {
		t.Error("Invalid number of cookies provided: ", session.HTTPClient.Jar.Cookies(cookieURL))
	}

	for _, cookie := range session.HTTPClient.Jar.Cookies(cookieURL) {
		switch cookie.Name {
		case "one":
			if cookie.Value != "two" {
				t.Error("Cookie value is not valid", cookie)
			}
		case "two":
			if cookie.Value != "three" {
				t.Error("Cookie value is not valid", cookie)
			}
		case "three":
			if cookie.Value != "four" {
				t.Error("Cookie value is not valid", cookie)
			}
		default:
			t.Error("We should not have any other cookies: ", cookie)
		}
	}

}

func TestPutInvalidURLSession(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Put("%../dir/", nil); err == nil {
		t.Error("Some how the request was valid to make request ", err)
	}
}
