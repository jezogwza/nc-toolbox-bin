package appliance

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorageClient(t *testing.T) {
	var storageClient *StorageClient = &StorageClient{}

	err := storageClient.InitStorageClient()
	if err != nil {
		fmt.Printf("Error creating storage client: %v\n", err)
	}
	assert.True(t, storageClient.purearray != nil, "Have a storage client")
}
