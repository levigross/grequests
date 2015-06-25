package grequests

import (
	"testing"
)

func TestBasicPutRequest(t *testing.T) {
	resp := Put("http://httpbin.org/put", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

}

func TestBasicAsyncPutRequest(t *testing.T) {
	resp := <-PutAsync("http://httpbin.org/put", nil)

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

}

func TestBasicPutUploadRequest(t *testing.T) {
	fd, err := FileUploadFromDisk("test_files/mypassword")

	if err != nil {
		t.Error("Unable to open file: ", err)
	}

	resp := <-PutAsync("http://httpbin.org/put",
		&RequestOptions{
			File: fd,
			Data: map[string]string{"One": "Two"},
		})

	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

}

func TestBasicPutUploadRequestInvalidURL(t *testing.T) {
	fd, err := FileUploadFromDisk("test_files/mypassword")

	if err != nil {
		t.Error("Unable to open file: ", err)
	}

	resp := Put("%../dir/",
		&RequestOptions{
			File: fd,
			Data: map[string]string{"One": "Two"},
		})

	if resp.Error == nil {
		t.Fatal("Somehow able to make the request")
	}
}
