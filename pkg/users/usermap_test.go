package users

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {

	var um UserMap
	err := um.LoadUsers("test/userlist.txt")
	if err != nil {
		fmt.Printf("Error is %s\n", err.Error())
	}
	//fmt.Println("Loaded users", um)
	for key := range um {
		assert.True(t, um[key].UserInfo.Name == key, "Should have be the same value as the key")
		assert.True(t, um[key].UserInfo.Password != "", "Should have a password value")
		assert.True(t, um[key].UserInfo.Role != "", "Should have a role value")
	}
}
func TestInvalidUsers(t *testing.T) {

	var um UserMap
	err := um.LoadUsers("test/baduserlist.txt")
	if err != nil {
		fmt.Printf("Error is %s\n", err.Error())
		assert.True(t, err != nil && err.Error() == "Invalid user admin. Used by a system account", "Should have failed to load the bad user list")
	}
}

func TestBadUsersRole(t *testing.T) {

	var um UserMap
	err := um.LoadUsers("test/baduserrole.txt")
	if err != nil {
		//fmt.Printf("Error is %s\n", err.Error())
		assert.True(t, err != nil && err.Error() == "Invalid role for user usera", "Should have failed to load the bad user role list")

	}
}

func TestListUsers(t *testing.T) {

	var um UserMap
	err := um.LoadUsers("test/userlist.txt")
	if err != nil {
		//fmt.Printf("Error is %s\n", err.Error())
		assert.False(t, err != nil && err.Error() == "Unable to list users", "Should have not errored,  failed to list users in list")
	}
	ulist := um.ListUsers()
	for _, u := range ulist {
		fmt.Println("User", u)
	}

}
