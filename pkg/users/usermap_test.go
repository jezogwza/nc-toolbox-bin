package users

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {

	var um UserMap
	um.Init()
	err := um.LoadUsers("test/userlist.txt")
	if err != nil {
		fmt.Printf("Error is %s\n", err.Error())
	}
	//fmt.Println("Loaded users", um)
	for key, _ := range um {
		assert.True(t, um[key].user.Name == key, "Should have be the same value as the key")
		assert.True(t, um[key].user.Password != "", "Should have a password value")
		assert.True(t, um[key].user.Role != "", "Should have a role value")
	}
}
func TestInvalidUsers(t *testing.T) {

	var um UserMap
	um.Init()
	err := um.LoadUsers("test/baduserlist.txt")
	if err != nil {
		//fmt.Printf("Error is %s\n", err.Error())
		assert.True(t, err != nil && err.Error() == "Invalid user admin. Used by a system account", "Should have failed to load the bad user list")
	}
}

func TestBadUsersRole(t *testing.T) {

	var um UserMap
	um.Init()
	err := um.LoadUsers("test/baduserrole.txt")
	if err != nil {
		//fmt.Printf("Error is %s\n", err.Error())
		assert.True(t, err != nil && err.Error() == "Invalid role for user usera", "Should have failed to load the bad user role list")

	}
}
