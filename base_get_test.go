package grequests

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

type BasicGetResponse struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
		Dnt            string `json:"Dst"`
		Host           string `json:"Host"`
		UserAgent      string `json:"User-Agent"`
		Hello          string `json:"Hello"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseNewHeader struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept           string `json:"Accept"`
		AcceptEncoding   string `json:"Accept-Encoding"`
		AcceptLanguage   string `json:"Accept-Language"`
		Dnt              string `json:"Dst"`
		Host             string `json:"Host"`
		UserAgent        string `json:"User-Agent"`
		XWonderfulHeader string `json:"X-Wonderful-Header"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseBasicAuth struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
		Dnt            string `json:"Dst"`
		Host           string `json:"Host"`
		UserAgent      string `json:"User-Agent"`
		Authorization  string `json:"Authorization"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseArgs struct {
	Args struct {
		Goodbye string `json:"Goodbye"`
		Hello   string `json:"Hello"`
	} `json:"args"`
	Headers struct {
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
		Dnt            string `json:"Dst"`
		Host           string `json:"Host"`
		UserAgent      string `json:"User-Agent"`
		Authorization  string `json:"Authorization"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type GetXMLSample struct {
	XMLName xml.Name `xml:"slideshow"`
	Title   string   `xml:"title,attr"`
	Date    string   `xml:"date,attr"`
	Author  string   `xml:"author,attr"`
	Slide   []struct {
		Type  string `xml:"type,attr"`
		Title string `xml:"title"`
	} `xml:"slide"`
}

type TestJSONCookies struct {
	Cookies struct {
		AnotherCookie string `json:"AnotherCookie"`
		TestCookie    string `json:"TestCookie"`
	} `json:"cookies"`
}

type MassiveJSONBlob struct {
	Type     string `json:"type"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			MAPBLKLOT string      `json:"MAPBLKLOT"`
			BLKLOT    string      `json:"BLKLOT"`
			BLOCKNUM  string      `json:"BLOCK_NUM"`
			LOTNUM    string      `json:"LOT_NUM"`
			FROMST    string      `json:"FROM_ST"`
			TOST      string      `json:"TO_ST"`
			STREET    string      `json:"STREET"`
			STTYPE    interface{} `json:"ST_TYPE"`
			ODDEVEN   string      `json:"ODD_EVEN"`
		} `json:"properties"`
		Geometry struct {
			Type        string `json:"type"`
			Coordinates []struct {
				Num0  []float64 `json:"0,omitempty"`
				Num1  []float64 `json:"1,omitempty"`
				Num2  []float64 `json:"2,omitempty"`
				Num3  []float64 `json:"3,omitempty"`
				Num4  []float64 `json:"4,omitempty"`
				Num5  []float64 `json:"5,omitempty"`
				Num6  []float64 `json:"6,omitempty"`
				Num7  []float64 `json:"7,omitempty"`
				Num8  []float64 `json:"8,omitempty"`
				Num9  []float64 `json:"9,omitempty"`
				Num10 []float64 `json:"10,omitempty"`
			} `json:"-"`
		} `json:"geometry"`
	} `json:"features"`
}

type GithubSelfJSON struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	Private          bool        `json:"private"`
	HTMLURL          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ForksURL         string      `json:"forks_url"`
	KeysURL          string      `json:"keys_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	TeamsURL         string      `json:"teams_url"`
	HooksURL         string      `json:"hooks_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	EventsURL        string      `json:"events_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BranchesURL      string      `json:"branches_url"`
	TagsURL          string      `json:"tags_url"`
	BlobsURL         string      `json:"blobs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	TreesURL         string      `json:"trees_url"`
	StatusesURL      string      `json:"statuses_url"`
	LanguagesURL     string      `json:"languages_url"`
	StargazersURL    string      `json:"stargazers_url"`
	ContributorsURL  string      `json:"contributors_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	CommitsURL       string      `json:"commits_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	CommentsURL      string      `json:"comments_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	ContentsURL      string      `json:"contents_url"`
	CompareURL       string      `json:"compare_url"`
	MergesURL        string      `json:"merges_url"`
	ArchiveURL       string      `json:"archive_url"`
	DownloadsURL     string      `json:"downloads_url"`
	IssuesURL        string      `json:"issues_url"`
	PullsURL         string      `json:"pulls_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	LabelsURL        string      `json:"labels_url"`
	ReleasesURL      string      `json:"releases_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	PushedAt         time.Time   `json:"pushed_at"`
	GitURL           string      `json:"git_url"`
	SSHURL           string      `json:"ssh_url"`
	CloneURL         string      `json:"clone_url"`
	SvnURL           string      `json:"svn_url"`
	Homepage         string      `json:"homepage"`
	Size             int         `json:"size"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Language         string      `json:"language"`
	HasIssues        bool        `json:"has_issues"`
	HasDownloads     bool        `json:"has_downloads"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	ForksCount       int         `json:"forks_count"`
	MirrorURL        interface{} `json:"mirror_url"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	Forks            int         `json:"forks"`
	OpenIssues       int         `json:"open_issues"`
	Watchers         int         `json:"watchers"`
	DefaultBranch    string      `json:"default_branch"`
	NetworkCount     int         `json:"network_count"`
	SubscribersCount int         `json:"subscribers_count"`
}

func TestGetNoOptions(t *testing.T) {
	resp, _ := Get("http://httpbin.org/get", nil)
	verifyOkResponse(resp, t)
}

func TestGetRequestHook(t *testing.T) {
	addHelloWorld := func(req *http.Request) error {
		req.Header.Add("Hello", "World")
		return nil
	}
	resp, _ := Get("http://httpbin.org/get",
		&RequestOptions{BeforeRequest: addHelloWorld})
	j := verifyOkResponse(resp, t)
	if j.Headers.Hello != "World" {
		t.Error("Hook Function failed")
	}
}

func TestGetNoOptionsCustomClient(t *testing.T) {
	resp, _ := Get("http://httpbin.org/get",
		&RequestOptions{HTTPClient: http.DefaultClient})
	verifyOkResponse(resp, t)
}

func TestGetCustomTLSHandshakeTimeout(t *testing.T) {
	ro := &RequestOptions{TLSHandshakeTimeout: 10 * time.Millisecond}
	if _, err := Get("https://httpbin.org", ro); err == nil {
		t.Error("unexpected: successful TLS Handshake")
	}
}

func TestGetCustomDialTimeout(t *testing.T) {
	ro := &RequestOptions{DialTimeout: time.Nanosecond}
	if _, err := Get("http://httpbin.org", ro); err == nil {
		t.Error("unexpected: successful connection")
	}
}

func TestGetProxy(t *testing.T) {
	ch := make(chan string, 1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch <- "real server"
	}))

	defer ts.Close()

	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch <- "proxy for " + r.URL.String()
	}))

	defer proxy.Close()

	pu, err := url.Parse(proxy.URL)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := Head(ts.URL, &RequestOptions{Proxies: map[string]*url.URL{pu.Scheme: pu}})

	defer http.DefaultTransport.(*http.Transport).CloseIdleConnections()

	if err != nil {
		t.Error("Unable to make request: ", err)
	}

	if resp.Ok != true {
		t.Error("Response is not OK for some reason: ", resp.StatusCode)
	}

	got := <-ch
	want := "proxy for " + ts.URL + "/"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestGetSyncInvalidProxyScheme(t *testing.T) {
	resp, err := Get("http://httpbin.org/get", &RequestOptions{Proxies: map[string]*url.URL{"gopher": nil}})
	if err != nil {
		t.Error("Request failed: ", err)
	}

	verifyOkResponse(resp, t)
}

func TestGetSyncNoOptions(t *testing.T) {
	resp, err := Get("http://httpbin.org/get", nil)
	if err != nil {
		t.Error("Request failed: ", err)
	}

	verifyOkResponse(resp, t)
}

func TestGetNoOptionsGzip(t *testing.T) {
	resp, _ := Get("https://httpbin.org/gzip", nil)
	verifyOkResponse(resp, t)
}

func TestGetWithCookies(t *testing.T) {
	resp, err := Get("http://httpbin.org/cookies",
		&RequestOptions{
			Cookies: []*http.Cookie{
				{
					Name:     "TestCookie",
					Value:    "Random Value",
					HttpOnly: true,
					Secure:   false,
				}, {
					Name:     "AnotherCookie",
					Value:    "Some Value",
					HttpOnly: true,
					Secure:   false,
				},
			},
		})

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &TestJSONCookies{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Cannot serialize cookie JSON: ", err)
	}

	if myJSONStruct.Cookies.TestCookie != "Random Value" {
		t.Errorf("Cookie value not set properly: %#v", myJSONStruct)
	}

	if myJSONStruct.Cookies.AnotherCookie != "Some Value" {
		t.Errorf("Cookie value not set properly: %#v", myJSONStruct)
	}

}

func TestGetWithCookiesCustomCookieJar(t *testing.T) {
	cookieJar, _ := cookiejar.New(nil)
	resp, err := Get("http://httpbin.org/cookies",
		&RequestOptions{
			CookieJar: cookieJar,
			Cookies: []*http.Cookie{
				{
					Name:     "TestCookie",
					Value:    "Random Value",
					HttpOnly: true,
					Secure:   false,
				}, {
					Name:     "AnotherCookie",
					Value:    "Some Value",
					HttpOnly: true,
					Secure:   false,
				},
			},
		})

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &TestJSONCookies{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Cannot serialize cookie JSON: ", err)
	}

	if myJSONStruct.Cookies.TestCookie != "Random Value" {
		t.Errorf("Cookie value not set properly: %#v", myJSONStruct)
	}

	if myJSONStruct.Cookies.AnotherCookie != "Some Value" {
		t.Errorf("Cookie value not set properly: %#v", myJSONStruct)
	}

}

func TestGetSession(t *testing.T) {
	session := NewSession(nil)

	resp, err := session.Get("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"one": "two"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	resp, err = session.Get("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"two": "three"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	resp, err = session.Get("http://httpbin.org/cookies/set", &RequestOptions{Params: map[string]string{"three": "four"}})

	if err != nil {
		t.Fatal("Cannot set cookie: ", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	cookieURL, err := url.Parse("http://httpbin.org")
	if err != nil {
		t.Error("We (for some reason) cannot parse the cookie URL")
	}

	if len(session.HTTPClient.Jar.Cookies(cookieURL)) != 3 {
		t.Error("Invalid number of cookies provided: ", session.HTTPClient.Jar.Cookies(cookieURL))
	}

	for _, cookie := range session.HTTPClient.Jar.Cookies(cookieURL) {
		switch cookie.Name {
		case "one":
			if cookie.Value != "two" {
				t.Error("Cookie value is not valid", cookie)
			}
		case "two":
			if cookie.Value != "three" {
				t.Error("Cookie value is not valid", cookie)
			}
		case "three":
			if cookie.Value != "four" {
				t.Error("Cookie value is not valid", cookie)
			}
		default:
			t.Error("We should not have any other cookies: ", cookie)
		}
	}

	session.CloseIdleConnections()

}

//func TestGetNoOptionsDeflate(t *testing.T) {
//	verifyOkResponse(<-GetAsync("http://httpbin.org/deflate", nil), t)
//}

func xmlASCIIDecoder(charset string, input io.Reader) (io.Reader, error) {
	return input, nil
}

func TestGetInvalidURL(t *testing.T) {
	resp, err := Get("%../dir/", &RequestOptions{Params: map[string]string{"1": "2"}})

	if err == nil {
		t.Error("Some how the request was valid to make request", err)
	}

	resp.ClearInternalBuffer() // This will panic without our nil checks
}

func TestGetInvalidURLNoParams(t *testing.T) {
	_, err := Get("%../dir/", nil)

	if err == nil {
		t.Error("Some how the request was valid to make request", err)
	}
}

func TestGetInvalidURLSession(t *testing.T) {
	session := NewSession(nil)

	if _, err := session.Get("%../dir/", nil); err == nil {
		t.Error("Some how the request was valid to make request ", err)
	}
}

func TestGetXMLSerialize(t *testing.T) {
	resp, err := Get("http://httpbin.org/xml", nil)

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	userXML := &GetXMLSample{}

	if err := resp.XML(userXML, xmlASCIIDecoder); err != nil {
		t.Error("Unable to consume the response as XML: ", err)
	}

	if userXML.Title != "Sample Slide Show" {
		t.Errorf("Invalid XML serialization %#v", userXML)
	}

	if err := resp.XML(int(123), nil); err == nil {
		t.Error("Still able to consume XML from used response")
	}

}

func TestGetCustomUserAgent(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1"}
	resp, _ := Get("http://httpbin.org/get", ro)
	jsonResp := verifyOkResponse(resp, t)
	if jsonResp.Headers.UserAgent != "LeviBot 0.1" {
		t.Error("User agent header not properly set")
	}
}

func TestGetBasicAuth(t *testing.T) {
	ro := &RequestOptions{Auth: []string{"Levi", "Bot"}}
	resp, err := Get("http://httpbin.org/get", ro)
	// Not the usual JSON so copy and paste from below

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicGetResponseBasicAuth{}

	err = resp.JSON(myJSONStruct)
	if err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.Headers.Authorization != "Basic TGV2aTpCb3Q=" {
		t.Error("Unable to set HTTP basic auth", myJSONStruct.Headers)
	}

}

func TestGetCustomHeader(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1",
		Headers: map[string]string{"X-Wonderful-Header": "1"}}
	resp, err := Get("http://httpbin.org/get", ro)
	// Not the usual JSON so copy and paste from below

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicGetResponseNewHeader{}

	err = resp.JSON(myJSONStruct)
	if err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.Headers.XWonderfulHeader != "1" {
		t.Error("Unable to set custom HTTP header", myJSONStruct.Headers)
	}
}

func TestGetInvalidSSLCertNoVerify(t *testing.T) {
	ro := &RequestOptions{InsecureSkipVerify: true}
	for _, badSSL := range []string{
		"https://self-signed.badssl.com/",
		"https://expired.badssl.com/",
		"https://wrong.host.badssl.com/",
	} {
		resp, err := Get(badSSL, ro)
		if err != nil {
			t.Error("Unable to make request", err)
		}
		if resp.Ok != true {
			t.Error("Request did not return OK")
		}
	}

}

func TestGetInvalidSSLCertNoVerifyNoOptions(t *testing.T) {
	for _, badSSL := range []string{
		"https://self-signed.badssl.com/",
		"https://expired.badssl.com/",
		"https://wrong.host.badssl.com/",
	} {
		resp, err := Get(badSSL, nil)
		if err == nil {
			t.Error("Unable to make request", err)
		}

		if resp.Ok == true {
			t.Error("Request did not return OK")
		}
	}
}

func TestGetInvalidSSLCertNoCompression(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1", DisableCompression: true}
	resp, err := Get("https://self-signed.badssl.com/", ro)

	if err == nil {
		t.Error("SSL verification worked when it shouldn't of", err)
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

}

func TestGetInvalidSSLCertWithCompression(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1", DisableCompression: false}
	resp, err := Get("https://self-signed.badssl.com/", ro)

	if err == nil {
		t.Error("SSL verification worked when it shouldn't of", err)
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

}

func TestErrorResponseNOOP(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1", DisableCompression: false}
	resp, err := Get("https://self-signed.badssl.com/", ro)

	if err == nil {
		t.Error("SSL verification worked when it shouldn't of", err)
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

	myJSONStruct := &BasicGetResponseArgs{}

	if err := resp.JSON(myJSONStruct); err == nil {
		t.Error("Somehow Able to convert to JSON", err)
	}

	if resp.Bytes() != nil {
		t.Error("Somehow byte buffer is working now (bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("Somehow byte buffer is working now (bytes)", resp.String())
	}

	resp.ClearInternalBuffer()

	if resp.Bytes() != nil {
		t.Error("Somehow byte buffer is working now (bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("Somehow byte buffer is working now (bytes)", resp.String())
	}

	userXML := &GetXMLSample{}

	if err := resp.XML(userXML, xmlASCIIDecoder); err == nil {
		t.Errorf("Somehow to consume the response as XML: %#v", userXML)
	}

	fileName := "randomFile"

	if err := resp.DownloadToFile(fileName); err == nil {
		t.Error("Somehow able to download to file: ", err)
	}

	var buf [1]byte

	if written, err := resp.Read(buf[:]); written != -1 && err == nil {
		t.Error("Somehow we were able to read from our error response")
	}

}

func TestGetInvalidSSLCertNoCompressionNoVerify(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1", InsecureSkipVerify: true, DisableCompression: true}
	resp, err := Get("https://self-signed.badssl.com/", ro)

	if err != nil {
		t.Error("SSL verification worked when it shouldn't of", err)
	}

	if resp.Ok != true {
		t.Error("Request did return OK")
	}

}

func TestGetInvalidSSLCertWithCompressionNoVerify(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1", InsecureSkipVerify: true, DisableCompression: false}
	resp, err := Get("https://self-signed.badssl.com/", ro)

	if err != nil {
		t.Error("SSL verification worked when it shouldn't of", err)
	}

	if resp.Ok != true {
		t.Error("Request did return OK")
	}

}

func TestGetInvalidSSLCert(t *testing.T) {
	ro := &RequestOptions{UserAgent: "LeviBot 0.1"}
	resp, err := Get("https://self-signed.badssl.com/", ro)

	if err == nil {
		t.Error("SSL verification worked when it shouldn't of", err)
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

}

func TestGetBasicArgs(t *testing.T) {
	ro := &RequestOptions{
		Params: map[string]string{"Hello": "World"},
	}
	resp, _ := Get("http://httpbin.org/get?Goodbye=World", ro)

	verifyOkArgsResponse(resp, t)

}

func TestGetBasicArgsQueryStruct(t *testing.T) {
	ro := &RequestOptions{
		QueryStruct: struct {
			Hello string `url:"Hello"`
		}{
			"World",
		},
	}
	resp, _ := Get("http://httpbin.org/get?Goodbye=World", ro)

	verifyOkArgsResponse(resp, t)

}

func TestGetBasicArgsQueryStructErr(t *testing.T) {
	ro := &RequestOptions{
		QueryStruct: 5,
	}
	resp, err := Get("http://httpbin.org/get?Goodbye=World", ro)

	if err == nil {
		t.Error("URL Parsing should have failed")
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

}

func TestGetBasicArgsQueryStructUrlQueryErr(t *testing.T) {
	ro := &RequestOptions{
		QueryStruct: 5,
	}
	resp, err := Get("http://httpbin.org/get?Goodbye=World%zz", ro)

	if err == nil {
		t.Error("URL Parsing should have failed")
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

}

func TestGetBasicArgsQueryStructUrlErr(t *testing.T) {
	ro := &RequestOptions{
		QueryStruct: 5,
	}
	resp, err := Get("%", ro)

	if err == nil {
		t.Error("URL Parsing should have failed")
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

}

func TestGetBasicArgsErr(t *testing.T) {
	ro := &RequestOptions{
		Params: map[string]string{"Hello": "World"},
	}
	resp, err := Get("http://httpbin.org/get?Goodbye=%zzz", ro)

	if err == nil {
		t.Error("URL Parsing should have failed")
	}

	if resp.Ok == true {
		t.Error("Request did return OK")
	}

}

func TestGetBasicArgsParams(t *testing.T) {
	ro := &RequestOptions{
		Params: map[string]string{"Hello": "World", "Goodbye": "World"},
	}
	resp, _ := Get("http://httpbin.org/get", ro)

	verifyOkArgsResponse(resp, t)
}

func TestGetBasicArgsParamsOverwrite(t *testing.T) {
	ro := &RequestOptions{
		Params: map[string]string{"Hello": "World", "Goodbye": "World"},
	}

	resp, _ := Get("http://httpbin.org/get?Hello=Nothing", ro)

	verifyOkArgsResponse(resp, t)
}

func TestGetFileDownload(t *testing.T) {
	resp, err := Get("http://httpbin.org/get", nil)

	fileName := "randomFile"

	if err := resp.DownloadToFile(fileName); err != nil {
		t.Error("Unable to download to file: ", err)
	}

	if err := resp.DownloadToFile("."); err == nil {
		t.Error("Able to create file '.'")
	}

	fd, err := os.Open(fileName)
	defer fd.Close()
	defer os.Remove(fileName)

	if err != nil {
		t.Error("Unable to open file to verify content ", err)
	}

	jsonDecoder := json.NewDecoder(fd)

	myJSONStruct := &BasicGetResponse{}

	if err := jsonDecoder.Decode(myJSONStruct); err != nil {
		t.Error("Unable to cocerce file to JSON ", err)
	}

	if myJSONStruct.URL != "http://httpbin.org/get" {
		t.Error("For some reason the URL isn't the same", myJSONStruct.URL)
	}

	if myJSONStruct.Headers.Host != "httpbin.org" {
		t.Error("The host header is invalid")
	}

	if resp.Bytes() != nil {
		t.Error("JSON decoding did not fully consume the response stream (Bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("JSON decoding did not fully consume the response stream (String)", resp.String())
	}

	if resp.StatusCode != 200 {
		t.Error("Response returned a non-200 code")
	}

}

func TestJsonConsumedResponse(t *testing.T) {
	resp, err := Get("http://httpbin.org/get", nil)

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.Bytes() == nil {
		t.Error("Unable to coerce value to bytes", resp.Bytes())
	}

	resp.ClearInternalBuffer()

	if err := resp.JSON(struct{}{}); err == nil {
		t.Error("Struct should not be able to hold JSON: ")
	}
}

func TestDownloadConsumedResponse(t *testing.T) {
	resp, err := Get("http://httpbin.org/get", nil)

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.Bytes() == nil {
		t.Error("Unable to coerce value to bytes")
	}

	resp.ClearInternalBuffer()

	if err := resp.DownloadToFile("randomFile"); err == nil {
		t.Error("Still able to download file: ", err)
	}

	defer os.Remove("randomFile")
}

func TestGetBytes(t *testing.T) {
	resp, err := Get("http://httpbin.org/get", nil)

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.Bytes() == nil {
		t.Error("JSON decoding did not fully consume the response stream")
	}

	if bytes.Compare(resp.Bytes(), resp.Bytes()) != 0 {
		t.Error("Body bytes have not been cached", resp.Bytes())
	}
}

func TestGetBytesNoBuffer(t *testing.T) {
	resp, err := Get("http://httpbin.org/get", nil)

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.Bytes() == nil {
		t.Error("Cannot coerce HTTP response to bytes")
	}

	if bytes.Compare(resp.Bytes(), resp.Bytes()) != 0 {
		t.Error("Body bytes have not been cached", resp.Bytes())
	}

	if err := resp.DownloadToFile("randomFile"); err != nil {
		t.Error("Unable to download file: ", err)
	}

	defer os.Remove("randomFile")

	resp.ClearInternalBuffer()

	if resp.Bytes() != nil {
		t.Error("Internal Buffer not cleaned up")
	}
}

func TestGetString(t *testing.T) {
	resp, err := Get("http://httpbin.org/get", nil)

	if err != nil {
		t.Error("Unable to make request", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	if resp.String() == "" {
		t.Error("Response Stream not returned as string", resp.String())
	}

	if resp.String() != resp.String() {
		t.Error("Body string have not been cached", resp.String())
	}

	if err := resp.DownloadToFile("randomFile"); err != nil {
		t.Error("Unable to download file: ", err)
	}

	defer os.Remove("randomFile")

	resp.ClearInternalBuffer()

	if resp.String() != "" {
		t.Error("Internal Buffer not cleaned up")
	}

}

func TestGetRedirectHeaderCopy(t *testing.T) {
	srv := httptest.NewServer(http.DefaultServeMux)
	http.HandleFunc("/foo", func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Custom") == "" {
			http.Error(w, "no custom header found", http.StatusBadRequest)
			return
		}
	})
	resp, err := Get(srv.URL+"/foo", &RequestOptions{Headers: map[string]string{"X-Custom": "1"}})

	if err != nil {
		t.Error("Redirect request failed", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	srv.Close()

}

func TestGetRedirectSecretHeaderNoCopy(t *testing.T) {
	srv := httptest.NewServer(http.DefaultServeMux)
	http.HandleFunc("/sec", func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Custom") == "" {
			http.Error(w, "no custom header found", http.StatusBadRequest)
			return
		}
	})
	resp, err := Get(srv.URL+"/sec", &RequestOptions{
		Headers: map[string]string{"X-Custom": "1"}, SensitiveHTTPHeaders: map[string]struct{}{"X-Custom": {}},
	})

	if err != nil {
		t.Error("Redirect request failed", err)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	srv.Close()

}

func TestMassiveJSONFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping massive JSON file download because short was called")
	}
	resp, err := Get("https://raw.githubusercontent.com/levigross/sf-city-lots-json/master/citylots.json", nil)
	if err != nil {
		t.Error("Request to massive JSON blob failed", err)
	}

	myjson := &MassiveJSONBlob{}

	if err := resp.JSON(myjson); err != nil {
		t.Error("Unable to serialize massive JSON blob", err)
	}

	if myjson.Type != "FeatureCollection" {
		t.Error("JSON did not properly serialize")
	}
}

func TestGitHubSelfJSON(t *testing.T) {
	resp, err := Get("https://api.github.com/repos/levigross/grequests", nil)
	if err != nil {
		t.Error("Request to reddit JSON blob failed", err)
	}

	myjson := &GithubSelfJSON{}

	if err := resp.JSON(myjson); err != nil {
		t.Error("Unable to serialize reddit JSON blob", err)
	}
}

func TestUnlimitedRedirects(t *testing.T) {
	srv := httptest.NewServer(http.DefaultServeMux)
	http.HandleFunc("/bar", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/bar", http.StatusMovedPermanently)
	})

	resp, err := Get(srv.URL+"/bar", &RequestOptions{Headers: map[string]string{"X-Custom": "1"}})

	if err == nil {
		t.Error("Redirect limitation failed", err)
	}

	if resp.Ok == true {
		t.Error("Request did not returned")
	}

	srv.Close()
}

func TestAuthStripOnRedirect(t *testing.T) {
	t.SkipNow()
	srv := httptest.NewServer(http.DefaultServeMux)
	http.HandleFunc("/test/", func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "" {
			http.Error(w, "Found Auth: "+req.Header.Get("Authorization"), http.StatusInternalServerError)
			return
		}

		if req.Header.Get("WWW-Authenticate") != "" {
			http.Error(w, "Found Auth: "+req.Header.Get("WWW-Authenticate"), http.StatusInternalServerError)
			return
		}

		if req.Header.Get("Proxy-Authorization") != "" {
			http.Error(w, "Found Auth: "+req.Header.Get("Proxy-Authorization"), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, "OK")
	})

	resp, err := Get(srv.URL+"/test", &RequestOptions{
		Auth:    []string{"one ", "two"},
		Headers: map[string]string{"WWW-Authenticate": "foo", "Proxy-Authorization": "bar"},
	})

	if err != nil {
		t.Error("Request had creds inside", err)
	}

	if resp.Ok != true {
		t.Error("Request had creds inside", resp.StatusCode, resp.String())
	}

	srv.Close()
}

func TestNoRedirect(t *testing.T) {
	srv := httptest.NewServer(http.DefaultServeMux)
	http.HandleFunc("/3tester/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	})

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("cancel redirection")
		},
	}

	_, err := Get(srv.URL+"/3tester/", &RequestOptions{
		HTTPClient: client,
	})

	if err == nil {
		t.Error("Request passed when it was supposed to fail on redirect", err)
	}

	srv.Close()

}

func verifyOkArgsResponse(resp *Response, t *testing.T) *BasicGetResponseArgs {
	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicGetResponseArgs{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.Args.Goodbye != "World" && myJSONStruct.Args.Hello != "World" {
		t.Error("Args not properly set", myJSONStruct.Args)
	}

	if myJSONStruct.URL != "http://httpbin.org/get?Goodbye=World&Hello=World" {
		t.Error("Url is not properly constructed", myJSONStruct.URL)
	}

	if resp.Bytes() != nil {
		t.Error("JSON decoding did not fully consume the response stream (Bytes)", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("JSON decoding did not fully consume the response stream (String)", resp.String())
	}

	if resp.StatusCode != 200 {
		t.Error("Response returned a non-200 code")
	}

	return myJSONStruct
}

func TestGetCustomRequestTimeout(t *testing.T) {
	ro := &RequestOptions{RequestTimeout: 2 * time.Nanosecond}
	if _, err := Get("http://httpbin.org", ro); err == nil {
		t.Error("unexpected: successful connection")
	}
}

func TestGetCustomRequestTimeoutContext(t *testing.T) {
	derContext := context.Background()
	ctx, cancel := context.WithTimeout(derContext, time.Microsecond)
	ro := &RequestOptions{Context: ctx}
	cancel()
	if _, err := Get("http://httpbin.org", ro); err == nil {
		t.Error("unexpected: successful connection")
	}
}

// verifyResponse will verify the following conditions
// 1. The request didn't return an error
// 2. The response returned an OK (a status code within the 200 range)
// 3. The output can be coerced to JSON (this may change later)
// It should only be run when testing GET request to http://httpbin.org/get expecting JSON
func verifyOkResponse(resp *Response, t *testing.T) *BasicGetResponse {
	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJSONStruct := &BasicGetResponse{}

	if err := resp.JSON(myJSONStruct); err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJSONStruct.Headers.Host != "httpbin.org" {
		t.Error("The host header is invalid")
	}

	if resp.Bytes() != nil {
		t.Errorf("JSON decoding did not fully consume the response stream (Bytes) %#v", resp.Bytes())
	}

	if resp.String() != "" {
		t.Error("JSON decoding did not fully consume the response stream (String)", resp.String())
	}

	if resp.StatusCode != 200 {
		t.Error("Response returned a non-200 code")
	}

	return myJSONStruct
}
