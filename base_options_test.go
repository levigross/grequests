package grequests

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type OptionsSuite struct {
	suite.Suite
}

func (s *OptionsSuite) TestOPTIONSRequest() {
	srv := newOptionsServer()
	defer srv.Close()

	resp, err := Options(srv.URL)
	s.Require().NoError(err)
	s.True(resp.Ok)
	s.Equal("GET, POST, PUT, DELETE, PATCH, OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
}

func (s *OptionsSuite) TestOptionsSessionCookies() {
	srv := newCookieSetServer()
	defer srv.Close()

	session := NewSession(nil)
	_, err := session.Options(srv.URL+"?one=two", &RequestOptions{})
	s.Require().NoError(err)
	_, err = session.Options(srv.URL+"?two=three", &RequestOptions{})
	s.Require().NoError(err)
	_, err = session.Options(srv.URL+"?three=four", &RequestOptions{})
	s.Require().NoError(err)

	cookieURL, err := url.Parse(srv.URL)
	s.Require().NoError(err)
	s.Len(session.HTTPClient.Jar.Cookies(cookieURL), 3)
}

func (s *OptionsSuite) TestOptionsInvalidURLSession() {
	session := NewSession(nil)
	_, err := session.Options("%../dir/", nil)
	s.Error(err)
}

func TestOptionsSuite(t *testing.T) {
	suite.Run(t, new(OptionsSuite))
}
