package users

import (
	"fmt"
	"os/user"
)

/** Only allow users with sudo priviliges */
func HavePriviliges() bool {
	u, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user")
		return false
	}
	// @@TODO - SHould check if teh user has access to sudo
	if u.Username == "root" {
		return true
	}
	return false
}
