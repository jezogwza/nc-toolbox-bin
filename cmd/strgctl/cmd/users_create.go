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
var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "The create subcommand is used to create users on  storage appliances",
	Long: `Given a file with a list of users and their roles, the command reconciles the list against
	the appliance. With the follwoing behavior :
	- If user exists it isnt created, no changes
	- If user doesnt exists it is created and a  password is generate, and the infrmation is stored 
	in the provided keyvault. 
	An exmaple of the command :
	
	'<cmd> users create --file userlist.txt --keyvault mykeyvault'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("users create called")
		fileName, _ := cmd.Flags().GetString(FILE)
		keyVault, _ := cmd.Flags().GetString(KEYVAULT)
		err := create(fileName, keyVault)
		if err != nil {
			fmt.Println("Error creating users", err)
		}
	},
}

func init() {
	usersCmd.AddCommand(userCreateCmd)
	flags := userCreateCmd.PersistentFlags()
	flags.StringP(FILE, "f", "", "The file with the list of users and their roles")
	flags.StringP(KEYVAULT, "k", "", "The keyvault to store the user information")
}

/*
Given a file with a list of users and their roles, the command reconciles the list against
*/
func create(fileName string, keyVault string) error {

	if !user.HavePriviliges() {
		return fmt.Errorf("Unable to create users from appliance. You do not have the priviliges")
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

	um, err = sclient.CreateUsers(um)
	if err != nil {
		return err

	}

	err = sclient.GetUsers(um)
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
