package appliance

import (
	"fmt"

	"k8s.io/client-go/clientcmd"
	"k8s.io/client-go/rest"
)

const (
	StorageServiceName = "strgmgmt"
	StorageNamespace   = "nc-system"
	StorageSecretKey   = "default"
)

//______________________________________________________//

func NewStorageClient() (PureArray, error) {
	// Get the Kubeconfig
	//
	kubeconfigPath := KUBECONFIG // Set your kubeconfig path

	// Load kubeconfig from file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		return nil, err
	}

	kClient := newKubernetesClient(config * rest.Config)
	endpointIP, err := kClient.GetServiceClusterIp(StorageServiceName, StorageNamespace)
	if err != nil {
		return nil, err
	}

	strorageApplianceName, err := kClient.GetStorageApplianceName(StorageNamespace)
	if err != nil {
		return nil, err
	}

	username, secretName, err := kClient.GetStorageUserInfo(strorageApplianceName, StorageNamespace)
	if err != nil {
		return nil, err
	}
	password, err := kClient.GetSecretValue(secretName, StorageNamespace, "default")
	if err != nil {
		return nil, err
	}
	logger
	purearray, err = NewPureArrayWithCredentials(endpointIP, username, password, logger)
	if err != nil {
		return nil, err
	}

	return purearray, nil
}
