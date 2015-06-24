package grequests

import "io"

const (
	localUserAgent = "GRequests 0.1"
)

// XMLCharDecoder is a helper type that takes a stream of bytes (not encoded in
// UTF-8) and returns a reader that encodes the bytes into UTF-8. This is done
// because Go's XML library only supports XML encoded in UTF-8
type XMLCharDecoder func(charset string, input io.Reader) (io.Reader, error)
