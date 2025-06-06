package grequests

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// FileUpload is a struct that is used to specify the file that a User
// wishes to upload.
type FileUpload struct {
	// Filename is the name of the file that you wish to upload. We use this to guess the mimetype as well as pass it onto the server
	FileName string

	// FileContents is happy as long as you pass it a io.ReadCloser (which most file use anyways)
	FileContents io.ReadCloser

	// FieldName is form field name
	FieldName string

	// FileMime represents which mimetime should be sent along with the file.
	// When empty, defaults to application/octet-stream
	FileMime string
}

// FileUploadFromDisk allows you to create a FileUpload struct slice by just specifying a location on the disk
func FileUploadFromDisk(fileName string) ([]FileUpload, error) {
	fd, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	return []FileUpload{{FileContents: fd, FileName: fileName}}, nil

}

// FileUploadFromGlob allows you to create a FileUpload struct slice by just specifying a glob location on the disk
// this function will gloss over all errors in the files and only upload the files that don't return errors from the glob
func FileUploadFromGlob(fileSystemGlob string) ([]FileUpload, error) {
	files, err := filepath.Glob(fileSystemGlob)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New("grequests: No files have been returned in the glob")
	}

	filesToUpload := make([]FileUpload, 0, len(files))

	for _, f := range files {
		s, err := os.Stat(f)
		if err != nil {
			// log.Printf("grequests: error stating file %s, skipping: %v", f, err) // Optional: log skipped files
			continue
		}
		if s.IsDir() {
			// log.Printf("grequests: %s is a directory, skipping", f) // Optional: log skipped directories
			continue
		}

		fd, err := os.Open(f)
		if err != nil {
			// log.Printf("grequests: error opening file %s for upload, skipping: %v", f, err) // Optional: log skipped files due to open error
			continue // Skip files that cannot be opened
		}
		filesToUpload = append(filesToUpload, FileUpload{FileContents: fd, FileName: filepath.Base(fd.Name())})

	}

	return filesToUpload, nil

}
