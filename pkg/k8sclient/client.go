package k8sclient

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	StorageApplianceResourceName = "storageappliance"
	StorageApplianceApiVersion   = "platform.afo-nc.microsoft.com/v1"
	StorageApplianceGroup        = "storageappliances.platform.afo-nc.microsoft.com"
	StorageApplianceNamespace    = "nc-system"
)

type KubernetesClient struct {
	kClient  kubernetes.Interface
	kdClient dynamic.Interface
	logger   logr.Logger
}

func newKubernetesClient(config *rest.Config) *KubernetesClient {
	client := kubernetes.NewForConfigOrDie(config)
	dynamicClient := dynamic.NewForConfigOrDie(config)
	// Create a logger using slogr
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log := zapr.NewLogger(zapLog)
	return &KubernetesClient{
		kClient:  client,
		kdClient: dynamicClient,
		logger:   log,
	}
}

func (k *KubernetesClient) GetSecretValue(secretName string, namespace string, secretKey string) (string, error) {

	secret, err := k.kClient.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Error retrieving Secret %s: %v\n", secretName, err)
		return "", err
	}

	for key, value := range secret.Data {
		if key == secretKey {
			return string(value), nil
		}
	}
	return "", fmt.Errorf("Secret %s does not contain key %s", secretName, secretKey)
}

func (k *KubernetesClient) GetServiceClusterIp(serviceName string, namespace string) (string, error) {

	secret, err := k.kClient.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Error retrieving Service %s: %v\n", serviceName, err)
		return "", err
	}

	return secret.Spec.ClusterIP, nil

}

/*
@@TODO: Uncomment this function once the custom resource is created
*/
func (k *KubernetesClient) GetStorageApplianceName(namespace string) (string, error) {
	/*
		customResourceList, err := k.kClientCmd.Dynamic().
			Resource(metav1.GroupVersionResource{
				Version:  StorageApplianceApiVersion,
				Resource: StorageApplianceResourceName,
			}).Namespace(StorageApplianceNamespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing custom resource: %v error: %v\n", StorageApplianceResourceName, err)
			return "", err
		}

		// Print custom resource names
		for _, cr := range customResourceList.Items {
			fmt.Printf("Custom resource name: %s\n", cr.GetName())
			return cr.GetName(), nil
		}
		return "", fmt.Errorf("No custom resources found")
	*/
	return "b37m25purestor1-v2mmktdbb5", nil
}

/*
@@TODO: Uncomment this function once the custom resource is created
*/
func (k *KubernetesClient) GetStorageUserInfo(resourceName string, namespace string) (string, string, error) {
	/*
		// Define the custom resource (e.g., MyCustomResource)
		customResourceName := resourceName
		customResourceNamespace := namespace

		fmt.Printf("Getting custom resource: %v from %v \n", customResourceName, customResourceNamespace)
		// Get the custom resource
		customResource, err := k.kdClient.
			Resource(schema.GroupVersionResource{
				Group:    StorageApplianceGroup,
				Version:  StorageApplianceApiVersion,
				Resource: StorageApplianceResourceName,
			}).
			Namespace(customResourceNamespace).
			Get(context.TODO(), customResourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error getting custom resource: %v\n", err)
			return "", "", err
		}

		// Extract relevant data from the custom resource
		// For example, if your custom resource has a field called "data", you can access it like this:
		data := customResource.Object["spec"]
		fmt.Printf("Custom resource spec: %v\n", data)

		return "", "", nil
	*/
	return "pureuser", "b37m25purestor1-v2mmktdbb5", nil
}
