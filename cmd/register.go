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

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new Account",
	Long: `Register a new Account to the agenda, the usage of it
you should specify the username, password, email and telephone
`,
	Run: func(cmd *cobra.Command, args []string) {
		debugLog("[command] register -u %s -p ***** -e %s -t %s ", username, email, telephone)

		// 检查非法参数
		validArgs := true
		infos := []string{"username", "password", "email", "telephone"}
		registerArgs := []string{username, password, email, telephone}
		isvalid := []func(string) bool{isvalidUsername, isvalidPassword, isvalidEmail, isvalidTelephone}
		for i, info := range infos {
			if !isvalid[i](registerArgs[i]) {
				validArgs = false
				fmt.Printf("[fail] the Field %s is invalid\n", info)
				debugLog("[fail] the Field %s is invalid\n", info)
			}
		}
		if !validArgs {
			return
		}

		// 检查是否重名
		var users entity.UserList
		entity.GetUsers(&users)
		for _, user := range users {
			if user.UserName == username {
				fmt.Printf("[fail] there has been account named %s\n", username)
				debugLog("[fail] there has been account named %s\n", username)
				return
			}
		}
		users = append(users, entity.User{
			UserName:  username,
			Password:  password,
			Email:     email,
			Telephone: telephone,
		})
		entity.SetUsers(&users)
		fmt.Printf("[success] new account %s has been added\n", username)
		debugLog("[success] new account %s has been added\n", username)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVarP(&username, "username", "u", "", "the username of new account")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "the password of new account")
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "the email of new account")
	registerCmd.Flags().StringVarP(&telephone, "telephone", "t", "", "the telephone of new account")
}

func isvalidUsername(username string) bool {
	return username != ""
}
func isvalidPassword(password string) bool {
	return password != ""
}
func isvalidEmail(email string) bool {
	return email != ""
}
func isvalidTelephone(telephone string) bool {
	return telephone != ""
}
