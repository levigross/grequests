package grequests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FileUploadSuite struct {
	suite.Suite
}

func (s *FileUploadSuite) TestErrorOpenFile() {
	fd, err := FileUploadFromDisk("I am Not A File")
	if err == nil {
		s.Fail("We are not getting an error back from our non existent file: ")
	}

	if fd != nil {
		s.Fail("We actually got back a pointer: ", fd)
	}
}

func (s *FileUploadSuite) TestGLOBFiles() {
	fd, err := FileUploadFromGlob("testdata/*")

	if err != nil {
		s.Fail("Got an invalid GLOB: %v", err)
	}

	if len(fd) != 2 {
		s.Fail("Some how we have more than two files in our glob %v %v", len(fd), fd)
	}
}

func (s *FileUploadSuite) TestInvalidGlob() {
	if _, err := FileUploadFromGlob("[-]"); err == nil {
		s.Fail("Somehow the glob worked")
	}
}

func (s *FileUploadSuite) TestNoGlobFiles() {
	if _, err := FileUploadFromGlob("notapath"); err == nil {
		s.Fail("Somehow got a valid GLOB")
	}
}

func (s *FileUploadSuite) TestGlobWithDir() {
	fd, err := FileUploadFromGlob("*test*")

	if err != nil {
		s.Fail("Glob failed %v", err)
	}

	for _, f := range fd {
		if f.FileName == "test_files" {
			s.Fail(fmt.Sprintf("%v is a dir (which cannot be uploaded)", f))
		}
	}

}

func TestFileUploadSuite(t *testing.T) {
	suite.Run(t, new(FileUploadSuite))
}
