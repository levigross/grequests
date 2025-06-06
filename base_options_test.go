package grequests

import "testing"

func TestBasicOPTIONSRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	resp, err := Options(httpbinURL + "/get")
	if err != nil {
		t.Fatalf("options failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}
