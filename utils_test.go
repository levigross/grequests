package grequests

import "testing"

func TestIsIdempotentMethod(t *testing.T) {
	if IsIdempotentMethod("POST") != false {
		t.Error("POST is not idempotent")
	}

	if IsIdempotentMethod("GET") != true {
		t.Error("GET is idempotent")
	}

	if IsIdempotentMethod("OPTIONS") != true {
		t.Error("OPTIONS is idempotent")
	}

	if IsIdempotentMethod("HEAD") != true {
		t.Error("HEAD is idempotent")
	}

}
