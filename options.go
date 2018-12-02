package grequests

// Option is the base type we use to apply request options
type Option interface {
	Apply(*RequestOptions)
}

type optionFunc func(*RequestOptions)

func (o optionFunc) Apply(r *RequestOptions) {
	o(r)
}

// FromRequestOptions is a function that you can use to upgrade your
// requests
func FromRequestOptions(r *RequestOptions) Option {
	return optionFunc(func(ro *RequestOptions) {
		*ro = *r
	})
}

// UserAgent sets the value of a requests user agent
func UserAgent(value string) Option {
	return optionFunc(func(ro *RequestOptions) {
		ro.UserAgent = value
	})
}
