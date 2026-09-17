package storage

import (
	"errors"
	"io"
)

// unconfiguredStorage lets the server boot and serve every existing
// endpoint normally when no image storage credentials are set, instead of
// failing to start. Only image upload/delete actually needs Storage, so
// those two admin/manager endpoints are the only things that fail - with a
// clear message - until real credentials are configured.
type unconfiguredStorage struct{}

func NewUnconfiguredStorage() Storage {
	return &unconfiguredStorage{}
}

var ErrStorageNotConfigured = errors.New(
	"image storage is not configured on this server",
)

func (s *unconfiguredStorage) Upload(io.Reader, string) (*UploadResult, error) {
	return nil, ErrStorageNotConfigured
}

func (s *unconfiguredStorage) Delete(string) error {
	return ErrStorageNotConfigured
}
