package grequests

import (
	"io"
	"os"
)

type FileUpload struct {
	FileName     string
	FileContents io.ReadCloser
}

// FileUploadFromDisk allows you to create a FileUpload struct by just specifying a location on the diskI
func FileUploadFromDisk(fileName string) (*FileUpload, error) {
	fd, err := os.Open(fileName)

	// I should really log the error
	if err != nil {
		return nil, err
	}

	return &FileUpload{FileContents: fd, FileName: fileName}, nil

}
