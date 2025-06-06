package grequests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddQueryStringParams(t *testing.T) {
	userURL, err := buildURLParams("https://www.google.com/", map[string]string{"1": "2", "3": "4"})

	if err != nil {
		assert.Fail(t, "URL Parse Error: ", err)
	}

	if userURL != "https://www.google.com/?1=2&3=4" {
		assert.Fail(t, "URL params not properly built", userURL)
	}
}

func TestSortAddQueryStringParams(t *testing.T) {
	userURL, err := buildURLParams("https://www.google.com/", map[string]string{"3": "4", "1": "2"})

	if err != nil {
		assert.Fail(t, "URL Parse Error: ", err)
	}

	if userURL != "https://www.google.com/?1=2&3=4" {
		assert.Fail(t, "URL params not properly built and sorted", userURL)
	}
}

func TestAddQueryStringParamsExistingParam(t *testing.T) {
	userURL, err := buildURLParams("https://www.google.com/?5=6", map[string]string{"3": "4", "1": "2"})

	if err != nil {
		assert.Fail(t, "URL Parse Error: ", err)
	}

	if userURL != "https://www.google.com/?1=2&3=4&5=6" {
		assert.Fail(t, "URL params not properly built and sorted", userURL)
	}
}
