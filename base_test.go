package grequests

import (
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

func TestGetNoOptions(t *testing.T) {
	resp := <-Get("http://httpbin.org/get", nil)
	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return return OK")
	}

	myJsonStruct := &BasicGetResponse{}

	err := resp.JSON(myJsonStruct)
	if err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJsonStruct.URL != "http://httpbin.org/get" {
		t.Error("For some reason the URL isn't the same", myJsonStruct.URL)
	}
}

func TestGetNoOptionsChannel(t *testing.T) {
	respChan := Get("http://httpbin.org/get", nil)
	select {
	case resp := <-respChan:
		if resp.Error != nil {
			t.Error("Unable to make request", resp.Error)
		}
		myJsonStruct := &BasicGetResponse{}

		err := resp.JSON(myJsonStruct)
		if err != nil {
			t.Error("Unable to coerce to JSON", err)
		}

		if myJsonStruct.URL != "http://httpbin.org/get" {
			t.Error("For some reason the URL isn't the same", myJsonStruct.URL)
		}
	}
}
