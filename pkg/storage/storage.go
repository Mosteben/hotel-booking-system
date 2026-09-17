package storage

import "io"

// UploadResult is the storage-provider-agnostic outcome of an upload.
// PublicID is whatever identifier the provider needs to later delete the
// file - callers must persist it if they want to be able to clean up.
type UploadResult struct {
	URL      string
	PublicID string
}

type Storage interface {
	Upload(file io.Reader, filename string) (*UploadResult, error)
	Delete(publicID string) error
}
