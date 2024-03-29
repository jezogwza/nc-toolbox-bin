/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage lifecycle actions on users on the storage appliance",
	Long: `Supports :
	- Creation of users, this includes password generation
	- Rotation of user passwords
	- Deletion of users`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("users called")
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
