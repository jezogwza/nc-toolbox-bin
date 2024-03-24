package appliance

const (
	StorageServiceName = "strgmgmt"
	StorageNamespace   = "nc-system"
	StorageSecretKey   = "default"

	StorageApplianceResourceName = "storageappliance"
	StorageApplianceApiVersion   = "platform.afo-nc.microsoft.com/v1"
	StorageApplianceGroup        = "storageappliances.platform.afo-nc.microsoft.com"
)
const KUBECONFIG string = "c:\\Users\\ropacheco\\.kube\\config"

const (
	StorageApplianceUserRoleArrayAdmin   string = "array_admin"
	StorageApplianceUserRoleStorageAdmin string = "storage_admin"
	StorageApplianceUserRoleReadOnly     string = "readonly"
	StorageApplianceUserRoleOpsAdmin     string = "ops_admin"
)

var StorageApplianceUserList = []string{"admin", "lma", "pureuser", "storage"}

// States of users

const (
	UserStateCreated      string = "created"
	UserStateDeleted      string = "deleted"
	UserStateNotExists    string = "not_exists"
	UserStatePasswdRotate string = "passwordrotated"
	UserStatePasswdStored string = "passwordstored"
	UserStateIntended     string = "intended"
)
