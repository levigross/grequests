package grequests

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type contextKey struct{}

type OptionFuncSuite struct{ suite.Suite }

func (s *OptionFuncSuite) TestOptions() {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{}
	ctx := context.WithValue(context.Background(), contextKey{}, "v")
	addr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1}
	proxyURL, _ := url.Parse("http://proxy")
	opts := []struct {
		opt   Option
		check func(*RequestOptions)
	}{
		{UserAgent("ua"), func(ro *RequestOptions) { s.Equal("ua", ro.UserAgent) }},
		{Files([]FileUpload{{FileName: "f", FileContents: io.NopCloser(strings.NewReader("data"))}}), func(ro *RequestOptions) { s.Len(ro.Files, 1); s.Equal("f", ro.Files[0].FileName) }},
		{JSON("j"), func(ro *RequestOptions) { s.Equal("j", ro.JSON) }},
		{XML("<x></x>"), func(ro *RequestOptions) { s.Equal("<x></x>", ro.XML) }},
		{DisableTLSCertValidation(), func(ro *RequestOptions) { s.True(ro.InsecureSkipVerify) }},
		{DisableCompression(), func(ro *RequestOptions) { s.True(ro.DisableCompression) }},
		{Host("h"), func(ro *RequestOptions) { s.Equal("h", ro.Host) }},
		{BasicAuth("u", "p"), func(ro *RequestOptions) { s.Equal([]string{"u", "p"}, ro.Auth) }},
		{IsAJAX(), func(ro *RequestOptions) { s.True(ro.IsAjax) }},
		{Cookies([]*http.Cookie{{Name: "n"}}), func(ro *RequestOptions) { s.Len(ro.Cookies, 1) }},
		{UseCookieJar(), func(ro *RequestOptions) { s.True(ro.UseCookieJar) }},
		{Proxies(map[string]*url.URL{"http": proxyURL}), func(ro *RequestOptions) { s.Equal(proxyURL, ro.Proxies["http"]) }},
		{TLSHandshakeTimeout(time.Second), func(ro *RequestOptions) { s.Equal(time.Second, ro.TLSHandshakeTimeout) }},
		{DialTimeout(time.Second), func(ro *RequestOptions) { s.Equal(time.Second, ro.DialTimeout) }},
		{DialKeepAlive(time.Second), func(ro *RequestOptions) { s.Equal(time.Second, ro.DialKeepAlive) }},
		{RequestTimeout(time.Second), func(ro *RequestOptions) { s.Equal(time.Second, ro.RequestTimeout) }},
		{HTTPClient(client), func(ro *RequestOptions) { s.Equal(client, ro.HTTPClient) }},
		{SensitiveHTTPHeaders("Foo"), func(ro *RequestOptions) { _, ok := ro.SensitiveHTTPHeaders["Foo"]; s.True(ok) }},
		{RedirectLimit(2), func(ro *RequestOptions) { s.Equal(2, ro.RedirectLimit) }},
		{RequestBody(strings.NewReader("b")), func(ro *RequestOptions) { s.NotNil(ro.RequestBody) }},
		{CookieJar(jar), func(ro *RequestOptions) { s.Equal(jar, ro.CookieJar) }},
		{Context(ctx), func(ro *RequestOptions) { s.Equal(ctx, ro.Context) }},
		{BeforeRequest(func(req *http.Request) error { return nil }), func(ro *RequestOptions) { s.NotNil(ro.BeforeRequest) }},
		{LocalAddr(addr), func(ro *RequestOptions) { s.Equal(addr, ro.LocalAddr) }},
	}
	for _, tc := range opts {
		ro := &RequestOptions{}
		tc.opt.Apply(ro)
		tc.check(ro)
	}
}

func (s *OptionFuncSuite) TestFromRequestOptions() {
	base := &RequestOptions{UserAgent: "ua"}
	opt := FromRequestOptions(base)
	ro := &RequestOptions{}
	opt.Apply(ro)
	s.Equal("ua", ro.UserAgent)
}

type ResponseMethodSuite struct{ suite.Suite }

func (s *ResponseMethodSuite) TestJSONAndBytes() {
	srv := newJSONServer(map[string]string{"foo": "bar"})
	defer srv.Close()
	resp, err := Get(context.Background(), srv.URL)
	s.Require().NoError(err)
	s.NotEmpty(resp.String())
	b := resp.Bytes()
	s.NotEmpty(b)
	var data map[string]string
	s.NoError(resp.JSON(&data))
	s.Equal("bar", data["foo"])
	resp.ClearInternalBuffer()
	s.Empty(resp.Bytes())
}

func (s *ResponseMethodSuite) TestXMLAndDownload() {
	xmlData := "<root><foo>bar</foo></root>"
	srv := newXMLServer(xmlData)
	defer srv.Close()
	resp, err := Get(context.Background(), srv.URL)
	s.Require().NoError(err)
	var v struct {
		Foo string `xml:"foo"`
	}
	s.NoError(resp.XML(&v, nil))
	s.Equal("bar", v.Foo)

	resp2, err := Get(context.Background(), srv.URL)
	s.Require().NoError(err)
	file := "test_download.tmp"
	s.NoError(resp2.DownloadToFile(file))
	defer func() { _ = os.Remove(file) }()
	contents, err := os.ReadFile(file)
	s.NoError(err)
	s.Equal(xmlData, string(contents))
}

type RedirectSuite struct{ suite.Suite }

func (s *RedirectSuite) TestAddRedirectFunctionality() {
	client := &http.Client{}
	ro := &RequestOptions{RedirectLimit: 2, SensitiveHTTPHeaders: map[string]struct{}{"Foo": {}}}
	addRedirectFunctionality(client, ro)
	req1 := httptest.NewRequest("GET", "http://x", nil)
	req1.Header.Set("Foo", "bar")
	req2 := httptest.NewRequest("GET", "http://y", nil)
	s.NoError(client.CheckRedirect(req2, []*http.Request{req1}))
	s.Empty(req2.Header.Get("Foo"))
	req3 := httptest.NewRequest("GET", "http://z", nil)
	err := client.CheckRedirect(req3, []*http.Request{req1, req2})
	s.Equal(ErrRedirectLimitExceeded, err)
}

type SessionSuite struct{ suite.Suite }

func (s *SessionSuite) TestCombineRequestOptions() {
	session := NewSession(&RequestOptions{UserAgent: "ua", Headers: map[string]string{"A": "1"}})
	ro := &RequestOptions{Headers: map[string]string{"B": "2"}}
	out := session.combineRequestOptions(ro)
	s.Equal("ua", out.UserAgent)
	s.Equal("2", out.Headers["B"])
	s.Equal("1", out.Headers["A"])
}

func TestExtraSuites(t *testing.T) {
	suite.Run(t, new(OptionFuncSuite))
	suite.Run(t, new(ResponseMethodSuite))
	suite.Run(t, new(RedirectSuite))
	suite.Run(t, new(SessionSuite))
	suite.Run(t, new(InternalFuncsSuite))
}

type InternalFuncsSuite struct{ suite.Suite }

func (s *InternalFuncsSuite) TestEscapeQuotes() {
	in := `"foo\\bar"`
	s.Equal(`\"foo\\\\bar\"`, escapeQuotes(in))
}

func (s *InternalFuncsSuite) TestCreateBasicJSONRequest() {
	ro := &RequestOptions{JSON: map[string]string{"foo": "bar"}}
	req, err := createBasicJSONRequest("POST", "http://x", ro)
	s.NoError(err)
	b, err := io.ReadAll(req.Body)
	s.NoError(err)
	s.JSONEq(`{"foo":"bar"}`, string(b))
	s.Equal("application/json", req.Header.Get("Content-Type"))
}

func (s *InternalFuncsSuite) TestCreateBasicXMLRequest() {
	xmlStr := "<root>ok</root>"
	ro := &RequestOptions{XML: xmlStr}
	req, err := createBasicXMLRequest("POST", "http://x", ro)
	s.NoError(err)
	b, err := io.ReadAll(req.Body)
	s.NoError(err)
	s.Equal(xmlStr, string(b))
	s.Equal("application/xml", req.Header.Get("Content-Type"))
}

func (s *InternalFuncsSuite) TestCreateFileUploadRequest() {
	ro := &RequestOptions{Files: []FileUpload{{FileName: ".txt", FileContents: io.NopCloser(strings.NewReader("data"))}}}
	req, err := createFileUploadRequest("PUT", "http://x", ro)
	s.NoError(err)
	b, err := io.ReadAll(req.Body)
	s.NoError(err)
	s.Equal("data", string(b))
	s.NotEmpty(req.Header.Get("Content-Type"))
}

func (s *InternalFuncsSuite) TestCreateMultiPartPostRequest() {
	ro := &RequestOptions{
		Files: []FileUpload{
			{FileName: "f1.txt", FileContents: io.NopCloser(strings.NewReader("one"))},
			{FileName: "f2.txt", FileContents: io.NopCloser(strings.NewReader("two")), FieldName: "custom"},
		},
		Data: map[string]string{"foo": "bar"},
	}
	req, err := createMultiPartPostRequest("POST", "http://x", ro)
	s.NoError(err)
	mr, err := req.MultipartReader()
	s.NoError(err)
	parts := map[string]string{}
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		s.NoError(err)
		data, err := io.ReadAll(p)
		s.NoError(err)
		parts[p.FormName()] = string(data)
	}
	s.Equal("one", parts["file1"])
	s.Equal("two", parts["custom"])
	s.Equal("bar", parts["foo"])
}

type errReader struct{ err error }

func (e errReader) Read(p []byte) (int, error) { return 0, e.err }
func (e errReader) Close() error               { return nil }

func (s *InternalFuncsSuite) TestCreateMultiPartPostRequestErrors() {
	_, err := createMultiPartPostRequest("POST", "http://x", &RequestOptions{Files: []FileUpload{{FileContents: nil}}})
	s.Error(err)

	_, err = createMultiPartPostRequest("POST", "http://x", &RequestOptions{Files: []FileUpload{{FileName: "f", FileContents: errReader{err: fmt.Errorf("bad")}}}})
	s.Error(err)
}

func (s *InternalFuncsSuite) TestCreateBasicJSONRequestError() {
	_, err := createBasicJSONRequest("POST", "http://x", &RequestOptions{JSON: make(chan int)})
	s.Error(err)
}

func (s *InternalFuncsSuite) TestCreateBasicXMLRequestError() {
	_, err := createBasicXMLRequest("POST", "http://x", &RequestOptions{XML: make(chan int)})
	s.Error(err)
}

func (s *InternalFuncsSuite) TestBuildRequestParamsAndQueryStruct() {
	q := make(chan string, 2)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q <- r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := buildRequest("GET", srv.URL, &RequestOptions{Params: map[string]string{"a": "b"}}, nil)
	s.NoError(err)
	s.Equal("a=b", <-q)

	type qs struct {
		A string `url:"a"`
	}
	_, err = buildRequest("GET", srv.URL, &RequestOptions{QueryStruct: qs{A: "1"}}, nil)
	s.NoError(err)
	s.Equal("a=1", <-q)
}

func (s *InternalFuncsSuite) TestResponseHelpers() {
	resp := &Response{Error: fmt.Errorf("e")}
	n, err := resp.Read(make([]byte, 1))
	s.Equal(-1, n)
	s.Error(err)
	s.EqualError(resp.DownloadToFile("x"), "e")
}

func (s *InternalFuncsSuite) TestProxySettings() {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	ro := RequestOptions{}
	u, err := ro.proxySettings(req)
	s.NoError(err)
	if u != nil {
		s.Contains(u.String(), "proxy")
	}

	proxyURL, _ := url.Parse("http://special")
	ro.Proxies = map[string]*url.URL{"http": proxyURL}
	u, err = ro.proxySettings(req)
	s.NoError(err)
	s.Equal(proxyURL, u)
}

func (s *InternalFuncsSuite) TestAddHTTPHeaders() {
	ro := &RequestOptions{Headers: map[string]string{"X": "Y"}, UserAgent: "ua", Host: "h", Auth: []string{"u", "p"}, IsAjax: true}
	req := httptest.NewRequest("GET", "http://x", nil)
	addHTTPHeaders(ro, req)
	s.Equal("Y", req.Header.Get("X"))
	s.Equal("ua", req.Header.Get("User-Agent"))
	s.Equal("h", req.Host)
	s.Equal("basic", strings.ToLower(req.Header.Get("Authorization")[:5]))
	s.Equal("XMLHttpRequest", req.Header.Get("X-Requested-With"))
}

func (s *InternalFuncsSuite) TestAddCookies() {
	ro := &RequestOptions{Cookies: []*http.Cookie{{Name: "n", Value: "v"}}}
	req := httptest.NewRequest("GET", "http://x", nil)
	addCookies(ro, req)
	s.Equal(1, len(req.Cookies()))
	s.Equal("n", req.Cookies()[0].Name)
}

func (s *InternalFuncsSuite) TestBuildURLStruct() {
	type qs struct {
		A string `url:"a"`
		B []int  `url:"b"`
	}
	out, err := buildURLStruct("http://x", qs{A: "1", B: []int{2, 3}})
	s.NoError(err)
	s.Equal("http://x?a=1&b=2&b=3", out)
}

func (s *InternalFuncsSuite) TestBuildHTTPClient() {
	ro := RequestOptions{}
	s.Equal(http.DefaultClient, BuildHTTPClient(ro))
	custom := &http.Client{Timeout: time.Second}
	s.Equal(custom, BuildHTTPClient(RequestOptions{HTTPClient: custom}))
	c := BuildHTTPClient(RequestOptions{UseCookieJar: true})
	s.NotEqual(http.DefaultClient, c)
	s.NotNil(c.Jar)
}

func (s *InternalFuncsSuite) TestCloseIdleConnections() {
	sess := NewSession(nil)
	sess.CloseIdleConnections()
}
