package models

import "os"

// Image common information about image.
type Image struct {
	File        *os.File
	Filename    string
	ContentType string
	Filesize    int64
}
