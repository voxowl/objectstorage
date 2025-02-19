package objectstorage

import "io"

type ListOpts struct {
	Limit int // 0 means no limit
}

// ObjectStorage is a generic interface for interacting with object storage services
type ObjectStorage interface {

	// Download an object by key
	Download(key string) (io.ReadCloser, error)

	// Upload an object
	Upload(key string, bytesReadSeeker io.Reader) error

	// List objects
	List(prefix string, opts ListOpts) ([]string, error)
}
