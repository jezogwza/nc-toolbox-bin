/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	FILE 		= "file"
	KEYVAULT 	= "keyvault"
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
		var err := create(fileName, keyVault)
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
func create(fileName string, keyVault string) (error) {
	/** Load the list of users from the file */
	var um UserMap
	um.init()
	var err = um.LoadUsers(fileName)
	if err != nil {
		return err
	}
}
