package grequests

import "io"

type FileUpload struct {
	FileName     string
	FileContents io.ReadCloser
}
