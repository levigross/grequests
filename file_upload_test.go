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

func TestGLOBFiles(t *testing.T) {
	fd, err := FileUploadFromGlob("testdata/*")

	if err != nil {
		t.Error("Got an invalid GLOB: ", err)
	}

	if len(fd) != 2 {
		t.Error("Some how we have more than two files in our glob", len(fd), fd)
	}
}

func TestInvalidGlob(t *testing.T) {
	if _, err := FileUploadFromGlob("[-]"); err == nil {
		t.Error("Somehow the glob worked")
	}
}

func TestNoGlobFiles(t *testing.T) {
	if _, err := FileUploadFromGlob("notapath"); err == nil {
		t.Error("Somehow got a valid GLOB")
	}
}

func TestGlobWithDir(t *testing.T) {
	fd, err := FileUploadFromGlob("*test*")

	if err != nil {
		t.Error("Glob failed", err)
	}

	for _, f := range fd {
		if f.FileName == "test_files" {
			t.Error(f, "is a dir (which cannot be uploaded)")
		}
	}

}
