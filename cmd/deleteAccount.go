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
	"strings"

	"github.com/MrFive5555/GO_agenda/entity"
	"github.com/spf13/cobra"
)

// deleteAccountCmd represents the deleteAccount command
var deleteAccountCmd = &cobra.Command{
	Use:   "deleteAccount",
	Short: "delete current account",
	Long: `delete the current accout, 
	and this account will be removed from all meetings that it participates,
	and all the meetings that are sponsored by this account will be deleted    `,
	Run: func(cmd *cobra.Command, args []string) {

		debugLog("[command] deleteAccount " + strings.Join(args, " "))

		// other argument
		if len(args) > 0 {
			fmt.Println("too many arguments")
			debugLog("too many arguments")
			return
		}

		var users entity.UserList
		entity.GetUsers(&users)

		var state entity.LogState
		entity.GetLogState(&state)

		if state.HasLogin == false {
			fmt.Println("[fail] you haven't logged in any account")
			debugLog("[fail] you haven't logged in any account")
			return
		}

		// 已登录的用户可以删除本用户账户（即销号）
		var newUsers entity.UserList
		for _, user := range users {
			if user.UserName == state.UserName {
				continue
			}
			newUsers = append(newUsers, entity.User{
				user.UserName,
				user.Password,
				user.Email,
				user.Telephone,
			})
		}
		entity.SetUsers(&newUsers)

		// 用户账户删除以后：
		// 以该用户为 发起者 的会议将被删除
		// 以该用户为 参与者 的会议将从 参与者 列表中移除该用户
		// 若因此造成会议 参与者 人数为0，则会议也将被删除。
		var meetings entity.MeetingList
		entity.GetMeeting(&meetings)

		toRemove := make([]bool, len(meetings))
		for index, meeting := range meetings {
			if state.UserName == meeting.Sponsors {
				// delete meeting
				toRemove[index] = true
			}
			participatorsList := strings.Split(meeting.Participators, ",")
			for _, participator := range participatorsList {
				if state.UserName == participator {
					if len(participatorsList) == 1 {
						// delete meeting
						toRemove[index] = true
					} else {
						// delete participator
						var newParticipators string
						for _, newParticipator := range participatorsList {
							if newParticipator != state.UserName {
								newParticipators += ","
								newParticipators += newParticipator
							}
						}
						// remove the first commas
						meetings[index].Participators = newParticipators[1:]
						entity.SetMeeting(&meetings)
						fmt.Printf("[success] delete participator %s from meeting %s\n", state.UserName, meeting.Title)
						debugLog("[success] delete participator " + state.UserName + " from meeting " + meeting.Title)
					}
				}
			}
		}

		deleteMeeting(toRemove)

		fmt.Println("[success] account deleted successfully")
		debugLog("[success] account deleted successfully")

		// 删除成功则退出系统登录状态。删除后，该用户账户不再存在。
		state.UserName = ""
		state.HasLogin = false
		entity.SetLogState(&state)

	},
}

func init() {
	rootCmd.AddCommand(deleteAccountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteAccountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteAccountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
