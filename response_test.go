package grequests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResponseSuite struct {
	suite.Suite
}

func (s *ResponseSuite) TestResponseOk() {
	statuses := []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}
	for _, status := range statuses {
		srv := newStatusServer(status)
		resp, err := Get(context.Background(), srv.URL)
		srv.Close()
		s.Require().NoError(err)
		s.Equal(status, resp.StatusCode)
		s.True(resp.Ok)
	}
}

func TestResponseSuite(t *testing.T) {
	suite.Run(t, new(ResponseSuite))
}
