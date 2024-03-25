/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	storage "github.com/jezogwza/nc-toolbox-bin/pkg/appliance"
	user "github.com/jezogwza/nc-toolbox-bin/pkg/users"
	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var userDeleteCmd = &cobra.Command{
	Use:   "create",
	Short: "The delete subcommand is used to delete users on  storage appliances",
	Long: `Given a file with a list of users and their roles, the command reconciles the list against
	the appliance. With the following behavior :
	- If user exists and matches the role, it deletes it
	- If user exists and doesnt match the role, it doesnt delete it
	- If user doesnt exists it is ignored
	- If the user has been deleted the Keyvault entries will be deleted.
	An exmaple of the command :
	
	'<cmd> users delete --file userlist.txt --keyvault mykeyvault'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("users delete called")
		fileName, _ := cmd.Flags().GetString(FILE)
		keyVault, _ := cmd.Flags().GetString(KEYVAULT)
		err := delete(fileName, keyVault)
		if err != nil {
			fmt.Println("Error deleting users", err)
		}
	},
}

func init() {
	usersCmd.AddCommand(userDeleteCmd)
	flags := userDeleteCmd.PersistentFlags()
	flags.StringP(FILE, "f", "", "The file with the list of users and their roles")
	flags.StringP(KEYVAULT, "k", "", "The keyvault where users information should be stored")
}

/*
Given a file with a list of users and their roles, the command reconciles the list against
*/
func delete(fileName string, keyVault string) error {

	if !user.HavePriviliges() {
		return fmt.Errorf("Unable to delete users from appliance. You do not have the priviliges")
	}
	/** Load the list of users from the file */
	fmt.Println("Loading Users")
	var um user.UserMap
	err := um.LoadUsers(fileName)
	if err != nil {
		return err
	}

	fmt.Println("Initializing Storage Client ")
	var sclient *storage.StorageClient = &storage.StorageClient{}
	err = sclient.InitClient()
	if err != nil {
		return err
	}

	fmt.Println("Deleting Users")
	um, err = sclient.DeleteUsers(um)
	if err != nil {
		return err
	}
	/*
		I should have users with passwords , that have not been
		delivered to keyvauts yet?
		 pureadmin setattr --password requires oldpassword

	*/

	/** Need to get the Storage Appliance client

	- Get Storage Appliance Client
	- CreateUers (Get UM as an input):
		Should be a wrapper
		Smart to know if user exists and only change the password in that case
		Should be able to create a user and set the password
	*/

	/* Once this is complete then
	um.StoreUsers(keyVault)
	*/
	return nil
}
