package grequests

import "testing"

func TestBasicPutRequest(t *testing.T) {
	httpbinURL, teardown := setupHttpbinServerTest(t)
	defer teardown()

	ro := &RequestOptions{Data: map[string]string{"key": "value"}}
	resp, err := Put(httpbinURL+"/put", FromRequestOptions(ro))
	if err != nil {
		t.Fatalf("put failed: %v", err)
	}
	if !resp.Ok {
		t.Fatalf("expected ok response, got %d", resp.StatusCode)
	}
}
