package users

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"regexp"

	labels "github.com/jezogwza/nc-toolbox-bin/pkg"
)

/*

GEt Group of the User

root@b37stg01c1mg01 [ /etc ]# grep superAccessGroup /etc/sudoers
%superAccessGroup ALL=(ALL) NOPASSWD: ALL
*/
/** Only allow users with sudo priviliges */
func HavePriviliges() bool {
	u, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user")
		return false
	}

	fmt.Println("HavePriviliges: user: %v uid: %d\n", u.Username, u.Uid)

	// Get the group IDs
	groupIDs, err := u.GroupIds()

	fmt.Println("HavePriviliges: found groups : %v \n", groupIDs)
	if err != nil {
		fmt.Println("Error getting groups the user belongs to:", err)
		return false
	}

	return checkIfGroupInSudoers(groupIDs)

}

func checkIfGroupInSudoers(groupids []string) bool {

	// Open the file
	fmt.Printf("Opening %s file:\n", labels.SUDOERS_FILE)
	file, err := os.Open(labels.SUDOERS_FILE)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer file.Close()

	fmt.Printf("Found  %v groups:\n", groupids)

	for _, gid := range groupids {
		// Define the pattern to search for
		pattern := regexp.MustCompile("%" + gid + " ALL=(ALL) NOPASSWD: ALL")

		// Read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("Found  %v line \n", line)
			// Check if the line contains the pattern
			if pattern.MatchString(line) {
				fmt.Println(line)
				return true
			}
		}

		// Check for errors during scanning
		if err := scanner.Err(); err != nil {
			fmt.Println("Error scanning file:", err)
		}
	}
	return false
}
