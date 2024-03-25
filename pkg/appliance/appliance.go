package appliance

import (
	"fmt"

	labels "github.com/jezogwza/nc-toolbox-bin/pkg"
	k8sclient "github.com/jezogwza/nc-toolbox-bin/pkg/k8sclient"
	umap "github.com/jezogwza/nc-toolbox-bin/pkg/users"
)

type Appliance interface {
	// CreateUsers creates local users on the array.
	// It is advised to create one user at a time unless client consuming this API
	// keeps track of successful creation of each user input.
	// If user exists already, it recreates it. This is to handle redeployment
	// There is no way to retrieve API token for a user without login as that user
	// One can retrieve his own API token, not others even if it has got admin role
	CreateUsers(umap.UserMap) (umap.UserMap, error)
	// GetUsers gets a list of all the current users fro mthe stroage array.
	GetUsers(umap.UserMap) error

	// DeleteUsera deletes the users in a UserMap.  Deleting a user that doesn't exist returns success
	DeleteUsers(umap.UserMap) (umap.UserMap, error)

	// DeleteUser deletes the given user.  Deleting a user that doesn't exist returns success
	DeleteUser(username string) error

	// ChangeUserPassword Change the password of a user. The newPassword cannot be empty.
	ChangeUserPassword(userName string, password string, newPassword string) error

	InitClient() error
}

type StorageClient struct {
	purearray *PureArray
}

func (sc *StorageClient) InitClient() error {
	// Get the Kubeconfig
	config, err := k8sclient.GetConfig() // Set your kubeconfig path
	if err != nil {
		fmt.Printf("Error retrieving  kubeconfig: %v\n", err)
		return err
	}

	kClient := k8sclient.NewKubernetesClient(config)
	endpointIP, err := kClient.GetServiceClusterIp(labels.StorageServiceName, labels.StorageNamespace)
	if err != nil {
		return err
	}

	strorageApplianceName, err := kClient.GetStorageApplianceName(labels.StorageNamespace)
	if err != nil {
		return err
	}

	username, secretName, err := kClient.GetStorageUserInfo(strorageApplianceName, labels.StorageNamespace)
	if err != nil {
		return err
	}
	password, err := kClient.GetSecretValue(secretName, labels.StorageNamespace, "default")
	if err != nil {
		return err
	}

	purearray, err := NewPureArrayWithCredentials(endpointIP, username, password, kClient.GetLogger())
	if err != nil {
		return err
	}

	sc.purearray = purearray
	return nil
}

// CreateUsers creates local users on the array.
func (sc *StorageClient) CreateUsers(um umap.UserMap) (umap.UserMap, error) {
	uList, err := sc.purearray.CreateUsers(um.ListUsers())
	if err != nil {
		return nil, err
	}
	um.PrepareUsers(uList)
	return um, nil
}

func (sc *StorageClient) DeleteUsers(um umap.UserMap) (umap.UserMap, error) {
	for key := range um {
		ctlUser := (um)[key]
		err := sc.purearray.DeleteUser(ctlUser.UserInfo.Name)
		if err != nil {
			return nil, err
		}
		ctlUser.UserState = labels.UserStateDeleted
	}

	return um, nil
}

// GetUsers gets a list of all the current users fro mthe stroage array.
func (sc *StorageClient) GetUsers(um umap.UserMap) error {
	// NEed to maap from the list of user to the UserMap to keep state and the relationship to keyvault
	uList, err := sc.purearray.GetUsers()
	if err != nil {
		return err
	}
	um.PrepareUsers(uList)
	return nil
}

// DeleteUser deletes the given user.  Deleting a user that doesn't exist returns success
func (sc *StorageClient) DeleteUser(username string) error {
	return sc.purearray.DeleteUser(username)
}

// ChangeUserPassword Change the password of a user. The newPassword cannot be empty.
func (sc *StorageClient) ChangeUserPassword(userName string, password string, newPassword string) error {
	return sc.purearray.ChangeUserPassword(userName, password, newPassword)
}
