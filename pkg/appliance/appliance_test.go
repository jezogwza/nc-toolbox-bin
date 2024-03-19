package appliance

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorageClient(t *testing.T) {
	StorageClient, err := NewStorageClient()
	if err != nil {
		fmt.Printf("Error creating storage client: %v\n", err)
	}
	assert.True(t, StorageClient != nil, "Have a storage client")
}
