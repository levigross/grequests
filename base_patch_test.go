package grequests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PatchSuite struct {
	suite.Suite
}

func (s *PatchSuite) TestPatchRequest() {
	srv := newPatchServer()
	defer srv.Close()

	resp, err := Patch(srv.URL, FromRequestOptions(&RequestOptions{Data: map[string]string{"one": "two"}}))
	s.Require().NoError(err)
	s.True(resp.Ok)
}

func (s *PatchSuite) TestPatchInvalidURLSession() {
	session := NewSession(nil)
	_, err := session.Patch("%../dir/", nil)
	s.Error(err)
}

func TestPatchSuite(t *testing.T) {
	suite.Run(t, new(PatchSuite))
}
