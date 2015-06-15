package grequests

import (
	"io"
	"os"
)

type FileUpload struct {
	FileName     string
	FileContents io.ReadCloser
}

// FileUploadFromDisk allows you to create a FileUpload struct by just specifying a location on the disk
// right now it ignores the error if it is unable to upload the file. This will change when I figure out
// how to better implement this API
func FileUploadFromDisk(fileName string) *FileUpload {
	fd, err := os.Open(fileName)

	// I should really log the error
	if err != nil {
		return nil
	}

	return &FileUpload{FileContents: fd, FileName: fileName}

}
