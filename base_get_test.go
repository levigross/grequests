package grequests

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GetSuite struct {
	suite.Suite
}

func (s *GetSuite) TestGetRequest() {
	srv := newGetServer()
	defer srv.Close()

	resp, err := Get(context.Background(), srv.URL)
	s.Require().NoError(err)
	s.True(resp.Ok)
}

func (s *GetSuite) TestGetInvalidURL() {
	_, err := Get(context.Background(), "%../dir/")
	s.Error(err)
}

func (s *GetSuite) TestGetSessionCookies() {
	srv := newCookieSetServer()
	defer srv.Close()

	session := NewSession(nil)
	_, err := session.Get(context.Background(), srv.URL+"?one=two", nil)
	s.Require().NoError(err)
	_, err = session.Get(context.Background(), srv.URL+"?two=three", nil)
	s.Require().NoError(err)

	cookieURL, err := url.Parse(srv.URL)
	s.Require().NoError(err)
	s.Len(session.HTTPClient.Jar.Cookies(cookieURL), 2)
}

func TestGetSuite(t *testing.T) {
	suite.Run(t, new(GetSuite))
}
