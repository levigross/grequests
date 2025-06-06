package grequests

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

type httpbinResponse struct {
	Args    map[string][]string `json:"args,omitempty"`
	Headers http.Header         `json:"headers,omitempty"`
	Origin  string              `json:"origin,omitempty"`
	URL     string              `json:"url,omitempty"`
	Data    string              `json:"data,omitempty"`
	Files   map[string]string   `json:"files,omitempty"`
	Form    map[string][]string `json:"form,omitempty"`
	JSON    interface{}         `json:"json,omitempty"`
}

type cookiesResponse struct {
	Cookies map[string]string `json:"cookies"`
}

type authResponse struct {
	Authenticated bool   `json:"authenticated"`
	User          string `json:"user,omitempty"`
}

func createHttpbinTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Helper function to write JSON response
	writeJSON := func(w http.ResponseWriter, data interface{}, statusCode int) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(data)
	}

	// /get handler
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		origin := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			origin = xff
		}

		resp := httpbinResponse{
			Args:    r.URL.Query(),
			Headers: r.Header,
			Origin:  origin,
			URL:     r.URL.String(),
		}

		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(gz).Encode(resp)
		} else if strings.Contains(acceptEncoding, "deflate") {
			w.Header().Set("Content-Encoding", "deflate")
			fl := flate.NewWriter(w, flate.DefaultCompression)
			defer fl.Close()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(fl).Encode(resp)
		} else {
			writeJSON(w, resp, http.StatusOK)
		}
	})

	// /post, /put, /patch handlers
	handlePostPutPatch := func(w http.ResponseWriter, r *http.Request) {
		origin := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			origin = xff
		}

		resp := httpbinResponse{
			Args:    r.URL.Query(),
			Headers: r.Header,
			Origin:  origin,
			URL:     r.URL.String(),
		}

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		resp.Data = string(body)

		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
			r.ParseForm()
			resp.Form = r.PostForm
		} else if strings.HasPrefix(contentType, "application/json") {
			var jsonData interface{}
			json.Unmarshal(body, &jsonData)
			resp.JSON = jsonData
		} else if strings.HasPrefix(contentType, "multipart/form-data") {
			r.ParseMultipartForm(32 << 20) // 32MB max memory
			if r.MultipartForm != nil {
				resp.Form = r.MultipartForm.Value
				resp.Files = make(map[string]string)
				for key, fhs := range r.MultipartForm.File {
					if len(fhs) > 0 {
						// For simplicity, just take the first file's content if multiple are uploaded with the same name
						file, err := fhs[0].Open()
						if err == nil {
							defer file.Close()
							fileBytes, _ := io.ReadAll(file)
							resp.Files[key] = string(fileBytes) // Storing as string, httpbin often shows content
						}
					}
				}
			}
		}
		writeJSON(w, resp, http.StatusOK)
	}
	mux.HandleFunc("/post", handlePostPutPatch)
	mux.HandleFunc("/put", handlePostPutPatch)
	mux.HandleFunc("/patch", handlePostPutPatch)

	// /delete handler
	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		origin := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			origin = xff
		}
		resp := httpbinResponse{
			Args:    r.URL.Query(),
			Headers: r.Header,
			Origin:  origin,
			URL:     r.URL.String(),
		}
		writeJSON(w, resp, http.StatusOK)
	})

	// /cookies handler
	mux.HandleFunc("/cookies", func(w http.ResponseWriter, r *http.Request) {
		cookies := make(map[string]string)
		for _, cookie := range r.Cookies() {
			cookies[cookie.Name] = cookie.Value
		}
		writeJSON(w, cookiesResponse{Cookies: cookies}, http.StatusOK)
	})

	// /cookies/set handler
	mux.HandleFunc("/cookies/set", func(w http.ResponseWriter, r *http.Request) {
		for name, values := range r.URL.Query() {
			if len(values) > 0 {
				http.SetCookie(w, &http.Cookie{Name: name, Value: values[0], Path: "/"})
			}
		}
		http.Redirect(w, r, "/cookies", http.StatusFound)
	})

	// /redirect/:n handler
	mux.HandleFunc("/redirect/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/redirect/"), "/")
		if len(parts) == 1 && parts[0] != "" {
			n, err := strconv.Atoi(parts[0])
			if err == nil {
				if n > 0 {
					http.Redirect(w, r, fmt.Sprintf("/redirect/%d", n-1), http.StatusFound)
					return
				}
				if n == 0 {
					http.Redirect(w, r, "/get", http.StatusFound)
					return
				}
			}
		}
		http.Error(w, "Invalid redirect parameter", http.StatusBadRequest)
	})

	// /basic-auth/:user/:passwd handler
	mux.HandleFunc("/basic-auth/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		// Expected: basic-auth, user, passwd
		if len(pathParts) < 3 {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		expectedUser := pathParts[1]
		expectedPasswd := pathParts[2]

		user, passwd, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
			return
		}

		if user == expectedUser && passwd == expectedPasswd {
			writeJSON(w, authResponse{Authenticated: true, User: user}, http.StatusOK)
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
		}
	})

	// /status/:code handler
	mux.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/status/"), "/")
		if len(parts) == 1 && parts[0] != "" {
			code, err := strconv.Atoi(parts[0])
			if err == nil && code >= 100 && code <= 599 {
				// For some codes like 204, 304, writing a body is not allowed.
				// http.ResponseWriter handles this automatically if we just call WriteHeader.
				if http.StatusText(code) == "" { // Invalid code not in standard library
					http.Error(w, fmt.Sprintf("Invalid status code: %d", code), http.StatusBadRequest)
					return
				}

				if code == http.StatusNoContent || code == http.StatusNotModified {
					w.WriteHeader(code)
				} else {
					w.WriteHeader(code)
					// httpbin returns the status code description in the body for some codes
					// For simplicity, we might not add body for all, or just a generic one.
					// Let's try to return what httpbin does for common cases.
					fmt.Fprintf(w, "%d %s", code, http.StatusText(code))

				}
				return
			}
		}
		http.Error(w, "Invalid status code parameter", http.StatusBadRequest)
	})

	// /headers handler
	mux.HandleFunc("/headers", func(w http.ResponseWriter, r *http.Request) {
		origin := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			origin = xff
		}
		resp := struct { // httpbin /headers has a slightly different structure
			Headers http.Header `json:"headers"`
			Origin  string      `json:"origin,omitempty"`
		}{
			Headers: r.Header,
			Origin:  origin,
		}
		writeJSON(w, resp, http.StatusOK)
	})

	// /user-agent handler
	mux.HandleFunc("/user-agent", func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			UserAgent string `json:"user-agent"`
		}{
			UserAgent: r.Header.Get("User-Agent"),
		}
		writeJSON(w, resp, http.StatusOK)
	})

	// /gzip handler
	mux.HandleFunc("/gzip", func(w http.ResponseWriter, r *http.Request) {
		origin := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			origin = xff
		}
		respData := httpbinResponse{
			Args:    r.URL.Query(),
			Headers: r.Header, // Note: httpbin /gzip includes original request headers
			Origin:  origin,
			URL:     r.URL.String(),
		}
		// The actual httpbin /gzip response also includes a "gzipped": true field.
		// Let's create a custom struct for this specific response.
		gzipResp := struct {
			httpbinResponse
			Gzipped bool `json:"gzipped"`
		}{
			httpbinResponse: respData,
			Gzipped:         true,
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json") // httpbin still serves json
		gz := gzip.NewWriter(w)
		defer gz.Close()
		json.NewEncoder(gz).Encode(gzipResp)
	})

	// /deflate handler
	mux.HandleFunc("/deflate", func(w http.ResponseWriter, r *http.Request) {
		origin := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			origin = xff
		}
		respData := httpbinResponse{
			Args:    r.URL.Query(),
			Headers: r.Header,
			Origin:  origin,
			URL:     r.URL.String(),
		}
		// Similar to gzip, httpbin /deflate has a "deflated": true field.
		deflateResp := struct {
			httpbinResponse
			Deflated bool `json:"deflated"`
		}{
			httpbinResponse: respData,
			Deflated:        true,
		}

		w.Header().Set("Content-Encoding", "deflate")
		w.Header().Set("Content-Type", "application/json")
		fl := flate.NewWriter(w, flate.DefaultCompression)
		defer fl.Close()
		json.NewEncoder(fl).Encode(deflateResp)
	})

	return httptest.NewServer(mux)
}

// main function for local testing if needed
/*
func main() {
	server := createHttpbinTestServer()
	defer server.Close()
	fmt.Println("Test server listening on:", server.URL)
	// Keep server running until interrupted
	select {}
}
*/
