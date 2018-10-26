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

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new Account",
	Long: `Register a new Account to the agenda, the usage of it
you should specify the username, password, email and telephone
`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查非法参数
		validArgs := true
		infos := []string{"username", "password", "email", "telephone"}
		register_args := []string{username, password, email, telephone}
		isvalid := []func(string) bool{isvalid_username, isvalid_password, isvalid_email, isvalid_telephone}
		for i, info := range infos {
			if !isvalid[i](register_args[i]) {
				validArgs = false
				fmt.Printf("[fail] the Field %s is invalid\n", info)
			}
		}
		if !validArgs {
			return
		}

		// 检查是否重名
		var users UserList
		GetUsers(&users)
		for _, user := range users {
			if user.UserName == username {
				fmt.Printf("[fail] there has been account named %s\n", username)
				return
			}
		}
		users = append(users, User{
			username,
			password,
			email,
			telephone,
		})
		SetUsers(&users)
		fmt.Printf("[sucess] new account %s has been added\n", username)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVarP(&username, "username", "u", "", "the username of new account")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "the password of new account")
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "the email of new account")
	registerCmd.Flags().StringVarP(&telephone, "telephone", "t", "", "the telephone of new account")
}

func isvalid_username(username string) bool {
	return username != ""
}
func isvalid_password(password string) bool {
	return password != ""
}
func isvalid_email(email string) bool {
	return email != ""
}
func isvalid_telephone(telephone string) bool {
	return telephone != ""
}
