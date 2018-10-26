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

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "log in to agenda",
	Long: `use your username and password to log in
to agenda`,
	Run: func(cmd *cobra.Command, args []string) {
		var users UserList
		GetUsers(&users)

		for _, user := range users {
			if user.UserName == username {
				if password == user.Password {
					var state LogState
					GetLogState(&state)
					if state.HasLogin {
						fmt.Printf("[fail] account (%s) has been logged in\n", state.UserName)
					} else {
						state.UserName = username
						state.HasLogin = true
						SetLogState(&state)
						fmt.Printf("[success] Log in account (%s)\n", username)
					}
				} else {
					fmt.Printf("[fail] (%s) password uncorrect\n", username)
				}
				return
			}
		}
		fmt.Printf(`[fail] there is not account (%s). you can use 
"agenda register -u %s -p [password] -e [email] -t [telephone]"
to register it
`, username, username)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "the username of your account")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "the password of your account")
}
