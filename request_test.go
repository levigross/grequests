package grequests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RequestSuite struct {
	suite.Suite
}

func (s *RequestSuite) TestAddQueryStringParams() {
	userURL, err := buildURLParams("https://www.google.com/", map[string]string{"1": "2", "3": "4"})
	if s.NoError(err) {
		s.Equal("https://www.google.com/?1=2&3=4", userURL)
	}
}

func (s *RequestSuite) TestSortAddQueryStringParams() {
	userURL, err := buildURLParams("https://www.google.com/", map[string]string{"3": "4", "1": "2"})
	if s.NoError(err) {
		s.Equal("https://www.google.com/?1=2&3=4", userURL)
	}
}

func (s *RequestSuite) TestAddQueryStringParamsExistingParam() {
	userURL, err := buildURLParams("https://www.google.com/?5=6", map[string]string{"3": "4", "1": "2"})
	if s.NoError(err) {
		s.Equal("https://www.google.com/?1=2&3=4&5=6", userURL)
	}
}

func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}
