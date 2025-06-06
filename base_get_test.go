package grequests

import (
	"testing"
)

// TestBasicGet verifies that a simple GET request returns a 200 status code
func TestBasicGet(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, err := Get(httpbinURL + "/get")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if !resp.Ok {
		t.Fatalf("expected ok response, got status %d", resp.StatusCode)
	}
}
