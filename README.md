# GRequests
A Go "clone" of the great and famous Requests library

[![Build Status](https://travis-ci.org/levigross/grequests.svg?branch=master)](https://travis-ci.org/levigross/grequests) [![GoDoc](https://godoc.org/github.com/levigross/grequests?status.svg)](https://godoc.org/github.com/levigross/grequests)
[![Coverage Status](https://coveralls.io/repos/levigross/grequests/badge.svg)](https://coveralls.io/r/levigross/grequests)
[![Join the chat at https://gitter.im/levigross/grequests](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/levigross/grequests?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

License
======

GRequests is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text

Features
========

- Responses can be serialized into JSON and XML
- Easy file uploads
- Easy file downloads
- Support for the following HTTP verbs `GET, HEAD, POST, PUT, DELETE, PATCH, OPTIONS`

Install
=======
`go get -u github.com/levigross/grequests`

Usage
======
`import "github.com/levigross/grequests"`

Basic Examples
=========
Basic GET request:

```go
resp, err := grequests.Get("http://httpbin.org/get", nil)
// You can modify the request by passing an optional RequestOptions struct

if err != nil {
	log.Fatalln("Unable to make request: ", err)
}

fmt.Println(resp.String())
// {
//   "args": {},
//   "headers": {
//     "Accept": "*/*",
//     "Host": "httpbin.org",
```

If an error occurs all of the other properties and methods of a `Response` will be `nil`

Quirks
=======
## Request Quirks

When passing parameters to be added to a URL, if the URL has existing parameters that *_contradict_* with what has been passed within `Params` – `Params` will be the "source of authority" and overwrite the contradicting URL parameter.

Lets see how it works...

```go
ro := &RequestOptions{
	Params: map[string]string{"Hello": "Goodbye"},
}
Get("http://httpbin.org/get?Hello=World", ro)
// The URL is now http://httpbin.org/get?Hello=Goodbye
```

## Response Quirks

Order matters! This is because `grequests.Response` is implemented as an `io.ReadCloser` which proxies the *http.Response.Body* `io.ReadCloser` interface. It also includes an internal buffer for use in `Response.String()` and `Response.Bytes()`.

Here are a list of methods that consume the *http.Response.Body* `io.ReadCloser` interface.

- Response.JSON
- Response.XML
- Response.DownloadToFile
- Response.Close
- Response.Read

The following methods make use of an internal byte buffer

- Response.String
- Response.Bytes

In the code below, once the file is downloaded – the `Response` struct no longer has access to the request bytes

```go
response := Get("http://some-wonderful-file.txt", nil)

if err := response.DownloadToFile("randomFile"); err != nil {
	log.Println("Unable to download file: ", err)
}

// At this point the .String and .Bytes method will return empty responses

response.Bytes() == nil // true
response.String() == "" // true

```

But if we were to call `response.Bytes()` or `response.String()` first, every operation will succeed until the internal buffer is cleared:

```go
response := Get("http://some-wonderful-file.txt", nil)

// This call to .Bytes caches the request bytes in an internal byte buffer – which can be used again and again until it is cleared
response.Bytes() == `file-bytes`
response.String() == "file-string"

// This will work because it will use the internal byte buffer
if err := resp.DownloadToFile("randomFile"); err != nil {
	log.Println("Unable to download file: ", err)
}

// Now if we clear the internal buffer....
response.ClearInternalBuffer()

// At this point the .String and .Bytes method will return empty responses

response.Bytes() == nil // true
response.String() == "" // true
```




# grequests
`import "github.com/levigross/grequests"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
Package grequests implements a friendly API over Go's existing net/http library




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func BuildHTTPClient(ro RequestOptions) *http.Client](#BuildHTTPClient)
* [func EnsureTransporterFinalized(httpTransport *http.Transport)](#EnsureTransporterFinalized)
* [func FileUploadFromDisk(fileName string) ([]FileUpload, error)](#FileUploadFromDisk)
* [func FileUploadFromGlob(fileSystemGlob string) ([]FileUpload, error)](#FileUploadFromGlob)
* [type FileUpload](#FileUpload)
* [type RequestOptions](#RequestOptions)
* [type Response](#Response)
  * [func Delete(url string, ro *RequestOptions) (*Response, error)](#Delete)
  * [func Get(url string, ro *RequestOptions) (*Response, error)](#Get)
  * [func Head(url string, ro *RequestOptions) (*Response, error)](#Head)
  * [func Options(url string, ro *RequestOptions) (*Response, error)](#Options)
  * [func Patch(url string, ro *RequestOptions) (*Response, error)](#Patch)
  * [func Post(url string, ro *RequestOptions) (*Response, error)](#Post)
  * [func Put(url string, ro *RequestOptions) (*Response, error)](#Put)
  * [func (r *Response) Bytes() []byte](#Response.Bytes)
  * [func (r *Response) ClearInternalBuffer()](#Response.ClearInternalBuffer)
  * [func (r *Response) Close() error](#Response.Close)
  * [func (r *Response) DownloadToFile(fileName string) error](#Response.DownloadToFile)
  * [func (r *Response) JSON(userStruct interface{}) error](#Response.JSON)
  * [func (r *Response) Read(p []byte) (n int, err error)](#Response.Read)
  * [func (r *Response) String() string](#Response.String)
  * [func (r *Response) XML(userStruct interface{}, charsetReader XMLCharDecoder) error](#Response.XML)
* [type Session](#Session)
  * [func NewSession(ro *RequestOptions) *Session](#NewSession)
  * [func (s *Session) CloseIdleConnections()](#Session.CloseIdleConnections)
  * [func (s *Session) Delete(url string, ro *RequestOptions) (*Response, error)](#Session.Delete)
  * [func (s *Session) Get(url string, ro *RequestOptions) (*Response, error)](#Session.Get)
  * [func (s *Session) Head(url string, ro *RequestOptions) (*Response, error)](#Session.Head)
  * [func (s *Session) Options(url string, ro *RequestOptions) (*Response, error)](#Session.Options)
  * [func (s *Session) Patch(url string, ro *RequestOptions) (*Response, error)](#Session.Patch)
  * [func (s *Session) Post(url string, ro *RequestOptions) (*Response, error)](#Session.Post)
  * [func (s *Session) Put(url string, ro *RequestOptions) (*Response, error)](#Session.Put)
* [type XMLCharDecoder](#XMLCharDecoder)

#### <a name="pkg-examples">Examples</a>
* [Package (AcceptInvalidTLSCert)](#example__acceptInvalidTLSCert)
* [Package (BasicAuth)](#example__basicAuth)
* [Package (BasicGet)](#example__basicGet)
* [Package (BasicGetCustomHTTPClient)](#example__basicGetCustomHTTPClient)
* [Package (Cookies)](#example__cookies)
* [Package (CustomHTTPHeader)](#example__customHTTPHeader)
* [Package (CustomUserAgent)](#example__customUserAgent)
* [Package (DownloadFile)](#example__downloadFile)
* [Package (PostFileUpload)](#example__postFileUpload)
* [Package (PostForm)](#example__postForm)
* [Package (PostJSONAJAX)](#example__postJSONAJAX)
* [Package (PostXML)](#example__postXML)
* [Package (Proxy)](#example__proxy)
* [Package (Session)](#example__session)
* [Package (UrlQueryParams)](#example__urlQueryParams)

#### <a name="pkg-files">Package files</a>
[base.go](/src/github.com/levigross/grequests/base.go) [file_upload.go](/src/github.com/levigross/grequests/file_upload.go) [request.go](/src/github.com/levigross/grequests/request.go) [response.go](/src/github.com/levigross/grequests/response.go) [session.go](/src/github.com/levigross/grequests/session.go) [utils.go](/src/github.com/levigross/grequests/utils.go) 



## <a name="pkg-variables">Variables</a>
``` go
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
        "Www-Authenticate":    {},
        "Authorization":       {},
        "Proxy-Authorization": {},
    }
)
```


## <a name="BuildHTTPClient">func</a> [BuildHTTPClient](/src/target/request.go?s=11478:11530#L405)
``` go
func BuildHTTPClient(ro RequestOptions) *http.Client
```
BuildHTTPClient is a function that will return a custom HTTP client based on the request options provided
the check is in UseDefaultClient



## <a name="EnsureTransporterFinalized">func</a> [EnsureTransporterFinalized](/src/target/utils.go?s=2482:2544#L79)
``` go
func EnsureTransporterFinalized(httpTransport *http.Transport)
```
EnsureTransporterFinalized will ensure that when the HTTP client is GCed
the runtime will close the idle connections (so that they won't leak)
this function was adopted from Hashicorp's go-cleanhttp package



## <a name="FileUploadFromDisk">func</a> [FileUploadFromDisk](/src/target/file_upload.go?s=625:687#L14)
``` go
func FileUploadFromDisk(fileName string) ([]FileUpload, error)
```
FileUploadFromDisk allows you to create a FileUpload struct slice by just specifying a location on the disk



## <a name="FileUploadFromGlob">func</a> [FileUploadFromGlob](/src/target/file_upload.go?s=1068:1136#L27)
``` go
func FileUploadFromGlob(fileSystemGlob string) ([]FileUpload, error)
```
FileUploadFromGlob allows you to create a FileUpload struct slice by just specifying a glob location on the disk
this function will gloss over all errors in the files and only upload the files that don't return errors from the glob




## <a name="FileUpload">type</a> [FileUpload](/src/target/file_upload.go?s=162:512#L2)
``` go
type FileUpload struct {
    // Filename is the name of the file that you wish to upload. We use this to guess the mimetype as well as pass it onto the server
    FileName string

    // FileContents is happy as long as you pass it a io.ReadCloser (which most file use anyways)
    FileContents io.ReadCloser

    // FieldName is form field name
    FieldName string
}
```
FileUpload is a struct that is used to specify the file that a User
wishes to upload.










## <a name="RequestOptions">type</a> [RequestOptions](/src/target/request.go?s=357:4258#L18)
``` go
type RequestOptions struct {

    // Data is a map of key values that will eventually convert into the
    // query string of a GET request or the body of a POST request.
    Data map[string]string

    // Params is a map of query strings that may be used within a GET request
    Params map[string]string

    // QueryStruct is a struct that encapsulates a set of URL query params
    // this paramter is mutually exclusive with `Params map[string]string` (they cannot be combined)
    // for more information please see https://godoc.org/github.com/google/go-querystring/query
    QueryStruct interface{}

    // Files is where you can include files to upload. The use of this data
    // structure is limited to POST requests
    Files []FileUpload

    // JSON can be used when you wish to send JSON within the request body
    JSON interface{}

    // XML can be used if you wish to send XML within the request body
    XML interface{}

    // Headers if you want to add custom HTTP headers to the request,
    // this is your friend
    Headers map[string]string

    // InsecureSkipVerify is a flag that specifies if we should validate the
    // server's TLS certificate. It should be noted that Go's TLS verify mechanism
    // doesn't validate if a certificate has been revoked
    InsecureSkipVerify bool

    // DisableCompression will disable gzip compression on requests
    DisableCompression bool

    // UserAgent allows you to set an arbitrary custom user agent
    UserAgent string

    // Host allows you to set an arbitrary custom host
    Host string

    // Auth allows you to specify a user name and password that you wish to
    // use when requesting the URL. It will use basic HTTP authentication
    // formatting the username and password in base64 the format is:
    // []string{username, password}
    Auth []string

    // IsAjax is a flag that can be set to make the request appear
    // to be generated by browser Javascript
    IsAjax bool

    // Cookies is an array of `http.Cookie` that allows you to attach
    // cookies to your request
    Cookies []*http.Cookie

    // UseCookieJar will create a custom HTTP client that will
    // process and store HTTP cookies when they are sent down
    UseCookieJar bool

    // Proxies is a map in the following format
    // *protocol* => proxy address e.g http => http://127.0.0.1:8080
    Proxies map[string]*url.URL

    // TLSHandshakeTimeout specifies the maximum amount of time waiting to
    // wait for a TLS handshake. Zero means no timeout.
    TLSHandshakeTimeout time.Duration

    // DialTimeout is the maximum amount of time a dial will wait for
    // a connect to complete.
    DialTimeout time.Duration

    // KeepAlive specifies the keep-alive period for an active
    // network connection. If zero, keep-alive are not enabled.
    DialKeepAlive time.Duration

    // RequestTimeout is the maximum amount of time a whole request(include dial / request / redirect)
    // will wait.
    RequestTimeout time.Duration

    // HTTPClient can be provided if you wish to supply a custom HTTP client
    // this is useful if you want to use an OAUTH client with your request.
    HTTPClient *http.Client

    // SensitiveHTTPHeaders is a map of sensitive HTTP headers that a user
    // doesn't want passed on a redirect.
    SensitiveHTTPHeaders map[string]struct{}

    // RedirectLimit is the acceptable amount of redirects that we should expect
    // before returning an error be default this is set to 30. You can change this
    // globally by modifying the `RedirectLimit` variable.
    RedirectLimit int

    // RequestBody allows you to put anything matching an `io.Reader` into the request
    // this option will take precedence over any other request option specified
    RequestBody io.Reader

    // CookieJar allows you to specify a special cookiejar to use with your request.
    // this option will take precedence over the `UseCookieJar` option above.
    CookieJar http.CookieJar

    // Context can be used to maintain state between requests https://golang.org/pkg/context/#Context
    Context context.Context
}
```
RequestOptions is the location that of where the data










## <a name="Response">type</a> [Response](/src/target/response.go?s=181:806#L4)
``` go
type Response struct {

    // Ok is a boolean flag that validates that the server returned a 2xx code
    Ok bool

    // This is the Go error flag – if something went wrong within the request, this flag will be set.
    Error error

    // We want to abstract (at least at the moment) the Go http.Response object away. So we are going to make use of it
    // internal but not give the user access
    RawResponse *http.Response

    // StatusCode is the HTTP Status Code returned by the HTTP Response. Taken from resp.StatusCode
    StatusCode int

    // Header is a net/http/Header structure
    Header http.Header
    // contains filtered or unexported fields
}
```
Response is what is returned to a user when they fire off a request







### <a name="Delete">func</a> [Delete](/src/target/base.go?s=1221:1283#L22)
``` go
func Delete(url string, ro *RequestOptions) (*Response, error)
```
Delete takes 2 parameters and returns a Response struct. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil


### <a name="Get">func</a> [Get](/src/target/base.go?s=300:359#L1)
``` go
func Get(url string, ro *RequestOptions) (*Response, error)
```
Get takes 2 parameters and returns a Response Struct. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil


### <a name="Head">func</a> [Head](/src/target/base.go?s=1841:1901#L38)
``` go
func Head(url string, ro *RequestOptions) (*Response, error)
```
Head takes 2 parameters and returns a Response channel. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil


### <a name="Options">func</a> [Options](/src/target/base.go?s=2151:2214#L46)
``` go
func Options(url string, ro *RequestOptions) (*Response, error)
```
Options takes 2 parameters and returns a Response struct. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil


### <a name="Patch">func</a> [Patch](/src/target/base.go?s=910:971#L14)
``` go
func Patch(url string, ro *RequestOptions) (*Response, error)
```
Patch takes 2 parameters and returns a Response struct. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil


### <a name="Post">func</a> [Post](/src/target/base.go?s=1533:1593#L30)
``` go
func Post(url string, ro *RequestOptions) (*Response, error)
```
Post takes 2 parameters and returns a Response channel. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil


### <a name="Put">func</a> [Put](/src/target/base.go?s=604:663#L6)
``` go
func Put(url string, ro *RequestOptions) (*Response, error)
```
Put takes 2 parameters and returns a Response struct. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil





### <a name="Response.Bytes">func</a> (\*Response) [Bytes](/src/target/response.go?s=4298:4331#L169)
``` go
func (r *Response) Bytes() []byte
```
Bytes returns the response as a byte array




### <a name="Response.ClearInternalBuffer">func</a> (\*Response) [ClearInternalBuffer](/src/target/response.go?s=4920:4960#L198)
``` go
func (r *Response) ClearInternalBuffer()
```
ClearInternalBuffer is a function that will clear the internal buffer that we use to hold the .String() and .Bytes()
data. Once you have used these functions – you may want to free up the memory.




### <a name="Response.Close">func</a> (\*Response) [Close](/src/target/response.go?s=1793:1825#L55)
``` go
func (r *Response) Close() error
```
Close is part of our ability to support io.ReadCloser if someone wants to make use of the raw body




### <a name="Response.DownloadToFile">func</a> (\*Response) [DownloadToFile](/src/target/response.go?s=2018:2074#L67)
``` go
func (r *Response) DownloadToFile(fileName string) error
```
DownloadToFile allows you to download the contents of the response to a file




### <a name="Response.JSON">func</a> (\*Response) [JSON](/src/target/response.go?s=3300:3353#L124)
``` go
func (r *Response) JSON(userStruct interface{}) error
```
JSON is a method that will populate a struct that is provided `userStruct` with the JSON returned within the
response body




### <a name="Response.Read">func</a> (\*Response) [Read](/src/target/response.go?s=1551:1603#L45)
``` go
func (r *Response) Read(p []byte) (n int, err error)
```
Read is part of our ability to support io.ReadCloser if someone wants to make use of the raw body




### <a name="Response.String">func</a> (\*Response) [String](/src/target/response.go?s=4568:4602#L186)
``` go
func (r *Response) String() string
```
String returns the response as a string




### <a name="Response.XML">func</a> (\*Response) [XML](/src/target/response.go?s=2792:2874#L101)
``` go
func (r *Response) XML(userStruct interface{}, charsetReader XMLCharDecoder) error
```
XML is a method that will populate a struct that is provided `userStruct` with the XML returned within the
response body




## <a name="Session">type</a> [Session](/src/target/session.go?s=125:314#L1)
``` go
type Session struct {
    // RequestOptions is global options
    RequestOptions *RequestOptions

    // HTTPClient is the client that we will use to request the resources
    HTTPClient *http.Client
}
```
Session allows a user to make use of persistent cookies in between
HTTP requests







### <a name="NewSession">func</a> [NewSession](/src/target/session.go?s=532:576#L8)
``` go
func NewSession(ro *RequestOptions) *Session
```
NewSession returns a session struct which enables can be used to maintain establish a persistent state with the
server
This function will set UseCookieJar to true as that is the purpose of using the session





### <a name="Session.CloseIdleConnections">func</a> (\*Session) [CloseIdleConnections](/src/target/session.go?s=4738:4778#L124)
``` go
func (s *Session) CloseIdleConnections()
```
CloseIdleConnections closes the idle connections that a session client may make use of




### <a name="Session.Delete">func</a> (\*Session) [Delete](/src/target/session.go?s=3120:3195#L88)
``` go
func (s *Session) Delete(url string, ro *RequestOptions) (*Response, error)
```
Delete takes 2 parameters and returns a Response struct. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil
A new session is created by calling NewSession with a request options struct




### <a name="Session.Get">func</a> (\*Session) [Get](/src/target/session.go?s=1776:1848#L58)
``` go
func (s *Session) Get(url string, ro *RequestOptions) (*Response, error)
```
Get takes 2 parameters and returns a Response Struct. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil
A new session is created by calling NewSession with a request options struct




### <a name="Session.Head">func</a> (\*Session) [Head](/src/target/session.go?s=4022:4095#L108)
``` go
func (s *Session) Head(url string, ro *RequestOptions) (*Response, error)
```
Head takes 2 parameters and returns a Response channel. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil
A new session is created by calling NewSession with a request options struct




### <a name="Session.Options">func</a> (\*Session) [Options](/src/target/session.go?s=4473:4549#L118)
``` go
func (s *Session) Options(url string, ro *RequestOptions) (*Response, error)
```
Options takes 2 parameters and returns a Response struct. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil
A new session is created by calling NewSession with a request options struct




### <a name="Session.Patch">func</a> (\*Session) [Patch](/src/target/session.go?s=2668:2742#L78)
``` go
func (s *Session) Patch(url string, ro *RequestOptions) (*Response, error)
```
Patch takes 2 parameters and returns a Response struct. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil
A new session is created by calling NewSession with a request options struct




### <a name="Session.Post">func</a> (\*Session) [Post](/src/target/session.go?s=3573:3646#L98)
``` go
func (s *Session) Post(url string, ro *RequestOptions) (*Response, error)
```
Post takes 2 parameters and returns a Response channel. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil
A new session is created by calling NewSession with a request options struct




### <a name="Session.Put">func</a> (\*Session) [Put](/src/target/session.go?s=2221:2293#L68)
``` go
func (s *Session) Put(url string, ro *RequestOptions) (*Response, error)
```
Put takes 2 parameters and returns a Response struct. These two options are:


	1. A URL
	2. A RequestOptions struct

If you do not intend to use the `RequestOptions` you can just pass nil
A new session is created by calling NewSession with a request options struct




## <a name="XMLCharDecoder">type</a> [XMLCharDecoder](/src/target/utils.go?s=1529:1605#L42)
``` go
type XMLCharDecoder func(charset string, input io.Reader) (io.Reader, error)
```
XMLCharDecoder is a helper type that takes a stream of bytes (not encoded in
UTF-8) and returns a reader that encodes the bytes into UTF-8. This is done
because Go's XML library only supports XML encoded in UTF-8














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
