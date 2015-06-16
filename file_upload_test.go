package grequests

import "testing"

func TestErrorOpenFile(t *testing.T) {
	fd, err := FileUploadFromDisk("I am Not A File")
	if err == nil {
		t.Error("We are not getting an error back from our non existent file: ")
	}

	if fd != nil {
		t.Error("We actually got back a pointer: ", fd)
	}
}
