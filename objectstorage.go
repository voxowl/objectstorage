package objectstorage

import "io"

type ListOpts struct {
	Limit *int32
}

// ObjectStorage is a generic interface for interacting with object storage services
type ObjectStorage interface {

	// Download an object by key
	Download(key string) ([]byte, error)

	// Upload an object
	Upload(key string, bytesReadSeeker io.ReadSeeker) error

	// List objects
	List(prefix string, opts ListOpts) ([]string, error)
}
