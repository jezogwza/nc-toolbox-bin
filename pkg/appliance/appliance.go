package appliance

import (
	"fmt"

	k8sclient "github.com/jezogwza/nc-toolbox-bin/pkg/k8sclient"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	StorageServiceName = "strgmgmt"
	StorageNamespace   = "nc-system"
	StorageSecretKey   = "default"
)

const KUBECONFIG string = "c:\\Users\\ropacheco\\.kube\\config"

//______________________________________________________//

func NewStorageClient() (*PureArray, error) {
	// Get the Kubeconfig
	//
	kubeconfigPath := KUBECONFIG // Set your kubeconfig path

	// Load kubeconfig from file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		return nil, err
	}

	kClient := k8sclient.NewKubernetesClient(config)
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

	purearray, err := NewPureArrayWithCredentials(endpointIP, username, password, kClient.GetLogger())
	if err != nil {
		return nil, err
	}

	return purearray, nil
}
