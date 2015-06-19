package grequests

import "io"

const (
	localUserAgent = "GRequests 0.1"
)

type XMLCharDecoder func(charset string, input io.Reader) (io.Reader, error)
