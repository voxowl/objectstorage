package digitalocean

import (
	"testing"

	"github.com/voxowl/objectstorage"
)

// TestDigitalOceanImplementsInterface is a compile-time check to ensure
// DigitalOceanObjectStorage implements the ObjectStorage interface
func TestDigitalOceanImplementsInterface(t *testing.T) {
	var _ objectstorage.ObjectStorage = (*DigitalOceanObjectStorage)(nil)
}
