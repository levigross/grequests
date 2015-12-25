package grequests

import (
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	localUserAgent = "GRequests 0.6"

	// Default value for net.Dialer Timeout
	dialTimeout = 30 * time.Second

	// Default value for request Timeout
	requestTimeout = 90 * time.Second

	// Default value for net.Dialer KeepAlive
	dialKeepAlive = 30 * time.Second

	// Default value for http.Transport TLSHandshakeTimeout
	tlsHandshakeTimeout = 10 * time.Second
)

var (
	// ErrRedirectLimitExceeded is the error returned when the request responded
	// with too many redirects
	ErrRedirectLimitExceeded = errors.New("grequests: Request exceeded redirect count")

	// RedirectLimit is a tunable variable that specifies how many times we can
	// redirect in response to a redirect. This is the global variable, if you
	// wish to set this on a request by request basis, set it within the
	// `RequestOptions` structure
	RedirectLimit = 30

	// SensitiveHTTPHeaders is a map of sensitive HTTP headers that a user
	// doesn't want passed on a redirect. This is the global variable, if you
	// wish to set this on a request by request basis, set it within the
	// `RequestOptions` structure
	SensitiveHTTPHeaders = map[string]struct{}{
		"WWW-Authenticate":    {},
		"Authorization":       {},
		"Proxy-Authorization": {},
	}
)

// XMLCharDecoder is a helper type that takes a stream of bytes (not encoded in
// UTF-8) and returns a reader that encodes the bytes into UTF-8. This is done
// because Go's XML library only supports XML encoded in UTF-8
type XMLCharDecoder func(charset string, input io.Reader) (io.Reader, error)

func addRedirectFunctionality(client *http.Client, ro *RequestOptions) {
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if ro.RedirectLimit == 0 {
			ro.RedirectLimit = RedirectLimit
		}

		if len(via) >= ro.RedirectLimit {
			return ErrRedirectLimitExceeded
		}

		if ro.SensitiveHTTPHeaders == nil {
			ro.SensitiveHTTPHeaders = SensitiveHTTPHeaders
		}

		for k, vv := range via[0].Header {
			// Is this a sensitive header?
			if _, found := ro.SensitiveHTTPHeaders[k]; found && !ro.RedirectLocationTrusted {
				continue
			}

			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}

		return nil
	}
}
