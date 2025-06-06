package grequests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorOpenFile(t *testing.T) {
	fd, err := FileUploadFromDisk("I am Not A File")
	if err == nil {
		assert.Fail(t, "We are not getting an error back from our non existent file: ")
	}

	if fd != nil {
		assert.Fail(t, "We actually got back a pointer: ", fd)
	}
}

func TestGLOBFiles(t *testing.T) {
	fd, err := FileUploadFromGlob("testdata/*")

	if err != nil {
		assert.Fail(t, "Got an invalid GLOB: ", err)
	}

	if len(fd) != 2 {
		assert.Fail(t, "Some how we have more than two files in our glob", len(fd), fd)
	}
}

func TestInvalidGlob(t *testing.T) {
	if _, err := FileUploadFromGlob("[-]"); err == nil {
		assert.Fail(t, "Somehow the glob worked")
	}
}

func TestNoGlobFiles(t *testing.T) {
	if _, err := FileUploadFromGlob("notapath"); err == nil {
		assert.Fail(t, "Somehow got a valid GLOB")
	}
}

func TestGlobWithDir(t *testing.T) {
	fd, err := FileUploadFromGlob("*test*")

	if err != nil {
		assert.Fail(t, "Glob failed", err)
	}

	for _, f := range fd {
		if f.FileName == "test_files" {
			assert.Fail(t, fmt.Sprintf("%v is a dir (which cannot be uploaded)", f))
		}
	}

}
