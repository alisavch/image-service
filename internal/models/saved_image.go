package models

import "os"

// SavedImage common information about image.
type SavedImage struct {
	File        *os.File
	Filename    string
	ContentType string
	Filesize    int64
	fileURL     string
}
