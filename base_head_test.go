package grequests

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HeadSuite struct {
	suite.Suite
}

func (s *HeadSuite) TestHeadRequest() {
	srv := newHeadServer()
	defer srv.Close()

	resp, err := Head(srv.URL)
	s.Require().NoError(err)
	s.True(resp.Ok)
	s.Equal("text/plain", resp.Header.Get("Content-Type"))
}

func (s *HeadSuite) TestHeadSessionCookies() {
	srv := newCookieSetServer()
	defer srv.Close()

	session := NewSession(nil)
	_, err := session.Head(srv.URL+"?one=two", &RequestOptions{})
	s.Require().NoError(err)
	_, err = session.Head(srv.URL+"?two=three", &RequestOptions{})
	s.Require().NoError(err)
	_, err = session.Head(srv.URL+"?three=four", &RequestOptions{})
	s.Require().NoError(err)

	cookieURL, err := url.Parse(srv.URL)
	s.Require().NoError(err)
	s.Len(session.HTTPClient.Jar.Cookies(cookieURL), 3)
}

func (s *HeadSuite) TestHeadInvalidURLSession() {
	session := NewSession(nil)
	_, err := session.Head("%../dir/", nil)
	s.Error(err)
}

func TestHeadSuite(t *testing.T) {
	suite.Run(t, new(HeadSuite))
}
