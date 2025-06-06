package grequests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PutSuite struct {
	suite.Suite
}

func (s *PutSuite) TestPutRequest() {
	srv := newPutServer()
	defer srv.Close()

	resp, err := Put(srv.URL, FromRequestOptions(&RequestOptions{Data: map[string]string{"one": "two"}}))
	s.Require().NoError(err)
	s.True(resp.Ok)
}

func (s *PutSuite) TestPutInvalidURLSession() {
	session := NewSession(nil)
	_, err := session.Put("%../dir/", nil)
	s.Error(err)
}

func TestPutSuite(t *testing.T) {
	suite.Run(t, new(PutSuite))
}
