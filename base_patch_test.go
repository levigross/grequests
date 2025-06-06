package grequests

import "testing"

func TestBasicPatchRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	ro := &RequestOptions{Data: map[string]string{"key": "value"}}
	resp, err := Patch(httpbinURL+"/patch", FromRequestOptions(ro))
	if err != nil {
		t.Fatalf("patch failed: %v", err)
	}
	if !resp.Ok {
		t.Fatalf("expected ok response, got %d", resp.StatusCode)
	}
}
