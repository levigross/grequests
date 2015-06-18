package grequests

import "io"

var (
	idempotentHTTPMethods = map[string]bool{"GET": true, "HEAD": true, "OPTIONS": true}
)

type XMLCharDecoder func(charset string, input io.Reader) (io.Reader, error)

func IsIdempotentMethod(httpMethod string) bool {
	_, ok := idempotentHTTPMethods[httpMethod]
	return ok
}
