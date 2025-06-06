package grequests

import "testing"

func TestBasicHeadRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, err := Head(httpbinURL + "/get")
	if err != nil {
		t.Fatalf("head failed: %v", err)
	}
	if !resp.Ok {
		t.Fatalf("expected ok response, got %d", resp.StatusCode)
	}
}
