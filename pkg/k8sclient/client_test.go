package k8sclient

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

const KUBECONFIG string = "/home/rodolfo/.kube/config"

func TestLoadConfig(t *testing.T) {

	kubeconfigPath := KUBECONFIG // Set your kubeconfig path

	// Load kubeconfig from file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}
	KubernetesClient := newKubernetesClient(config)
	assert.True(t, KubernetesClient != nil, "Have a k8s client")
}

func TestGetSecretValue(t *testing.T) {

	kubeconfigPath := KUBECONFIG // Set your kubeconfig path

	// Load kubeconfig from file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}
	KubernetesClient := newKubernetesClient(config)
	secretValue, err := KubernetesClient.GetSecretValue("b37m25purestor1-v2mmktdbb5", "nc-system", "default")
	fmt.Printf(" retrieving secret value: %v\n", secretValue)
	if err != nil {
		fmt.Printf("Error retrieving secret value: %v\n", err)
	}
	assert.True(t, secretValue == "FNdXeZJhlfBv", "Returned a secret value")
}

func TestGetServiceClusterIp(t *testing.T) {

	kubeconfigPath := KUBECONFIG // Set your kubeconfig path

	// Load kubeconfig from file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}
	KubernetesClient := newKubernetesClient(config)
	serviceValue, err := KubernetesClient.GetServiceClusterIp("strgmgmt", "nc-system")
	fmt.Printf(" retrieving service: %v\n", serviceValue)
	if err != nil {
		fmt.Printf("Error retrieving serviceValue value: %v\n", err)
	}
	assert.True(t, serviceValue == "10.96.0.32", "Returned a clusterIP value")
}

func TestGetStorageApplianceName(t *testing.T) {
	kubeconfigPath := KUBECONFIG // Set your kubeconfig path

	// Load kubeconfig from file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}
	KubernetesClient := newKubernetesClient(config)

	saname, err := KubernetesClient.GetStorageApplianceName("nc-system")
	fmt.Printf(" retrieving saname: %v\n", saname)
	if err != nil {
		fmt.Printf("Error retrieving storage appliance name : %v\n", err)
	}
	assert.True(t, saname == "b37m25purestor1-v2mmktdbb5", "Returned a storage appliance name")
}
func TestGetStorageUserInfo(t *testing.T) {

	kubeconfigPath := KUBECONFIG // Set your kubeconfig path

	// Load kubeconfig from file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}
	KubernetesClient := newKubernetesClient(config)
	username, resourcename, err := KubernetesClient.GetStorageUserInfo("b37m25purestor1-v2mmktdbb5", "nc-system")
	fmt.Printf(" retrieving resourcename: %v\n", resourcename)
	fmt.Printf(" retrieving username: %v\n", username)
	if err != nil {
		fmt.Printf("Error retrieving resource value: %v\n", err)
	}
	assert.True(t, username == "pureuser", "Returned the storage username")
	assert.True(t, resourcename == "b37m25purestor1-v2mmktdbb5", "Returned a storage secret value")
}
