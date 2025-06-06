package grequests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PostSuite struct {
	suite.Suite
}

func (s *PostSuite) TestPostRequest() {
	srv := newPostServer()
	defer srv.Close()

	resp, err := Post(context.Background(), srv.URL, FromRequestOptions(&RequestOptions{Data: map[string]string{"one": "two"}}))
	s.Require().NoError(err)
	s.True(resp.Ok)
}

func (s *PostSuite) TestPostSession() {
	srv := newPostServer()
	defer srv.Close()

	session := NewSession(nil)
	resp, err := session.Post(context.Background(), srv.URL, &RequestOptions{Data: map[string]string{"one": "two"}})
	s.Require().NoError(err)
	s.True(resp.Ok)
}

func (s *PostSuite) TestPostInvalidURL() {
	_, err := Post(context.Background(), "%../dir/")
	s.Error(err)
}

func TestPostSuite(t *testing.T) {
	suite.Run(t, new(PostSuite))
}
