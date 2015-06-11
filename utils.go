package grequests

var (
	idempotentHTTPMethods = map[string]bool{"GET": true, "HEAD": true, "OPTIONS": true}
)

func IsIdempotentMethod(httpMethod string) bool {
	_, ok := idempotentHTTPMethods[httpMethod]
	return ok
}
