// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/MrFive5555/GO_agenda/entity"
	"github.com/spf13/cobra"
)

// showallCmd represents the showall command
var showallCmd = &cobra.Command{
	Use:   "showall",
	Short: "show all the users",
	Long:  `show all the users who have registered before in agenda`,
	Run: func(cmd *cobra.Command, args []string) {

		debugLog("[command] showall")

		// other argument
		if len(args) > 0 {
			fmt.Println("too many arguments")
			debugLog("too many arguments")
			return
		}

		var users entity.UserList
		entity.GetUsers(&users)

		if len(users) <= 0 {
			fmt.Println("[success] There is no user in agenda, please register first")
			debugLog("[success] There is no user in agenda, please register first")
			return
		}
		for key, user := range users {
			fmt.Printf("[user %d]\n\tname:\t%s\n\temail:\t%s\n\ttel:\t%s\n", key+1, user.UserName, user.Email, user.Telephone)
		}
		fmt.Printf("[success] Done! A total of %d users have been shown\n", len(users))
		debugLog("[success] Done! A total of %d users have been shown\n", len(users))

	},
}

func init() {
	rootCmd.AddCommand(showallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
