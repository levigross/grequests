package grequests

import "testing"

// TestBasicPost verifies that a simple POST request with form data works
func TestBasicPost(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	opts := &RequestOptions{Data: map[string]string{"one": "two"}}
	resp, err := Post(httpbinURL+"/post", FromRequestOptions(opts))
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	if !resp.Ok {
		t.Fatalf("expected ok response, got status %d", resp.StatusCode)
	}
}
