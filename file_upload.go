package grequests

import (
	"io"
	"os"
)

// FileUpload is a struct that is used to specify the file that a User
// wishes to upload.
type FileUpload struct {
	// Filename is the name of the file that you wish to upload. We use this to guess the mimetype as well as pass it onto the server
	FileName string

	// FileContents is happy as long as you pass it a io.ReadCloser (which most file use anyways)
	FileContents io.ReadCloser
}

// FileUploadFromDisk allows you to create a FileUpload struct by just specifying a location on the diskI
func FileUploadFromDisk(fileName string) (*FileUpload, error) {
	fd, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	return &FileUpload{FileContents: fd, FileName: fileName}, nil

}
