package users

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"dev.azure.com/msazuredev/AzureForOperatorsIndustry/_git/nc-1p-core.git/services/credentialmanager/sdk"
)

/**

Name      Type   Role
admin     local  array_admin
lma       local  readonly
pureuser  local  array_admin
storage   local  storage_admin

pureuser@b37int1a1pu01> pureadmin create ropacheco --role readonly
Enter password:
Retype password:

*/

/* This stuff is in the controller so schoose to copy instead of improting */
const (
	StorageApplianceUserRoleArrayAdmin   string = "array_admin"
	StorageApplianceUserRoleStorageAdmin string = "storage_admin"
	StorageApplianceUserRoleReadOnly     string = "readonly"
	StorageApplianceUserRoleOpsAdmin     string = "ops_admin"
)

var StorageApplianceUserList = []string{"admin", "lma", "pureuser", "storage"}

type User struct {
	Name string

	// Generated by Pure client
	Password string

	// Allowed values: readonly, storage_admin, and array_admin.
	// ops_admin not supported yet
	Role string
	// Provided by the Pure client.
	ApiToken string
}

type CtlUser struct {
	user           User
	secretLocation string
}

type UserMap map[string]CtlUser

// Read the user name list from a file and generate a password for each user
func (um *UserMap) init() {
	*um = make(map[string]CtlUser)
}

// Read the user name list from a file and generate a password for each user
func (um *UserMap) LoadUsers(filename string) error {
	um.init()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		n := strings.Split(scanner.Text(), ",")
		err := um.validateUser(n)
		if err != nil {
			return err
		}
		u := CtlUser{
			user: User{
				Name:     n[0],
				Role:     (n[1]),
				Password: um.generatePasswd(),
			},
			secretLocation: "",
		}

		(*um)[n[0]] = u
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

/*
Ensure the list does not contain users that are not allowed
Such as users created by the Nexus Storage Appliance Controller
or teh Admin user of teh Storage Appliance
Enusre that the roles associated with the users are valid roles
*/
func (um *UserMap) validateUser(n []string) error {
	fmt.Printf("Validating user %s with role %s\n", n[0], n[1])
	if slices.Contains(StorageApplianceUserList, n[0]) {
		return fmt.Errorf("Invalid user %s. Used by a system account", n[0])
	}

	if (n[1]) != StorageApplianceUserRoleArrayAdmin &&
		(n[1]) != StorageApplianceUserRoleStorageAdmin &&
		(n[1]) != StorageApplianceUserRoleReadOnly &&
		(n[1]) != StorageApplianceUserRoleOpsAdmin {
		return fmt.Errorf("Invalid role for user %s", n[0])
	}

	return nil
}

/*
* Walk through the Users and Store their information in the
Keyvault.

*
*/
func (um *UserMap) StoreUsers(keyVault string) error {
	return nil
}

/*
*

*
 */
func (um *UserMap) ListUsers() []User {
	ulist := make([]User, len(*um))
	for key, _ := range *um {
		ulist = append(ulist, (*um)[key].user)
	}
	return ulist
}

func (um *UserMap) PrepareUsers(ulist []User) error {
	for key, _ := range *um {
		ctluser := (*um)[key]
		ctluser.secretLocation = "@TODO GENERATE THE KEYVAUKT ENTRY NAME FOR THIS USER"
	}
	return nil
}

// Generate a password that meets Security requirements
// and limitation swithin a storage appliance
func (um *UserMap) generatePasswd() string {
	passwordMinLength := int32(12)
	passwordMinUpperCase := int32(3)
	passwordMinNumeric := int32(2)
	passwordMinSpecialChar := int32(1)
	passwordSpecialCharList := "!@#$%^&*()_+"
	var passwordRequirements *sdk.PasswordRequirements = &sdk.PasswordRequirements{
		MinLength:       &passwordMinLength,
		MinUpperCase:    &passwordMinUpperCase,
		MinNumeric:      &passwordMinNumeric,
		MinSpecialChar:  &passwordMinSpecialChar,
		SpecialCharList: &passwordSpecialCharList,
	}
	var pwd, err = GeneratePassword(passwordRequirements)
	if err != nil {
		return "SomeGarbagePassword123!@#"
	}
	return pwd
}
