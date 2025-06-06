package grequests

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DeleteSuite struct {
	suite.Suite
}

func (s *DeleteSuite) TestDeleteRequest() {
	srv := newDeleteServer()
	defer srv.Close()

	resp, err := Delete(context.Background(), srv.URL)
	s.Require().NoError(err)
	s.True(resp.Ok)
}

func (s *DeleteSuite) TestDeleteSessionCookies() {
	srv := newCookieSetServer()
	defer srv.Close()

	session := NewSession(nil)
	_, err := session.Get(context.Background(), srv.URL+"?one=two", nil)
	s.Require().NoError(err)
	_, err = session.Get(context.Background(), srv.URL+"?two=three", nil)
	s.Require().NoError(err)
	_, err = session.Get(context.Background(), srv.URL+"?three=four", nil)
	s.Require().NoError(err)

	_, err = session.Delete(context.Background(), srv.URL, nil)
	s.Require().NoError(err)

	cookieURL, err := url.Parse(srv.URL)
	s.Require().NoError(err)
	s.Len(session.HTTPClient.Jar.Cookies(cookieURL), 3)
}

func (s *DeleteSuite) TestDeleteInvalidURLSession() {
	session := NewSession(nil)
	_, err := session.Delete(context.Background(), "%../dir/", nil)
	s.Error(err)
}

func TestDeleteSuite(t *testing.T) {
	suite.Run(t, new(DeleteSuite))
}
