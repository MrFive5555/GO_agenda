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

// removeParticipatorCmd represents the removeParticipator command
var removeParticipatorCmd = &cobra.Command{
	Use:   "removeParticipator",
	Short: "remove some participators to one of your meeting",
	Long: `To remove some participators to one of your meeting,
you should specify the title and the new participators `,
	Run: func(cmd *cobra.Command, args []string) {

		// 检查非法参数
		validArgs := true
		infos := []string{"title", "participators"}
		registerArgs := []string{title, participators}
		isvalid := []func(string) bool{isvalidTitle, isvalidParticipators}
		for i, info := range infos {
			if !isvalid[i](registerArgs[i]) {
				validArgs = false
				fmt.Printf("[fail] the Field %s is invalid\n", info)
				debugLog("[fail] the Field " + info + " is invalid")
			}
		}
		if !validArgs {
			return
		}

		// 获得当前登录用户信息
		var state entity.LogState
		var me string
		entity.GetLogState(&state)
		if state.HasLogin {
			me = state.UserName
		} else {
			fmt.Printf("[fail] you haven't loged in\n")
			debugLog("[fail] you haven't logged in")
			return
		}

		// 检查自己发起的主题为title的会议是否存在，并获取该会议信息
		var meetings entity.MeetingList
		entity.GetMeeting(&meetings)
		var myMeeting int
		validMeeting := false
		for i, m := range meetings {
			if (m.Title == title) && (m.Sponsors == me) {
				validMeeting = true
				myMeeting = i
				break
			}
		}
		if !validMeeting {
			fmt.Printf("[fail] %s does not have sponsored the meeting %s\n", me, title)
			debugLog("[fail] " + me + " does not have sponsored the meeting " + title)
			return
		}

		// 检查用户是否已注册
		var users entity.UserList
		entity.GetUsers(&users)
		participatorsList := strings.Split(participators, ",")
		for _, p := range participatorsList {
			validParticipator := false
			for _, user := range users {
				if user.UserName == p {
					validParticipator = true
					break
				}
			}
			if !validParticipator {
				fmt.Printf("[fail] participator %s does not exist\n", p)
				debugLog("[fail] participator " + p + " does not exist")
				return
			}
		}

		// 删除会议参与者，人数为0则删除会议

		toRemove := make([]bool, len(meetings))
		for _, participator := range participatorsList {
			meetingParticipators := meetings[myMeeting].Participators
			mpList := strings.Split(meetingParticipators, ",")
			isParticipated := false
			for _, mp := range mpList {
				if mp == participator {
					if len(mpList) == 1 {
						// delete meeting
						toRemove[myMeeting] = true
						fmt.Printf("[success] delete participator %s from meeting %s\n", mp, meetings[myMeeting].Title)
						debugLog("[success] delete participator " + mp + " from meeting " + meetings[myMeeting].Title)
						isParticipated = true
						break
					} else {
						// delete participator
						var newParticipators string
						for _, newParticipator := range mpList {
							if newParticipator != mp {
								newParticipators += ","
								newParticipators += newParticipator
							}
						}
						// remove the first commas
						meetings[myMeeting].Participators = newParticipators[1:]
						entity.SetMeeting(&meetings)
						fmt.Printf("[success] delete participator %s from meeting %s\n", mp, meetings[myMeeting].Title)
						debugLog("[success] delete participator " + mp + " from meeting " + meetings[myMeeting].Title)
						isParticipated = true
						break
					}
				}
			}
			if !isParticipated {
				fmt.Printf("[fail] participator %s is not in the meeting %s\n", participator, meetings[myMeeting].Title)
				debugLog("[fail] participator " + participator + " is not in the meeting " + meetings[myMeeting].Title)
			}
		}

		deleteMeeting(toRemove)

	},
}

func init() {
	rootCmd.AddCommand(removeParticipatorCmd)

	removeParticipatorCmd.Flags().StringVarP(&title, "title", "t", "", "the title of new meeting")
	removeParticipatorCmd.Flags().StringVarP(&participators, "participators", "p", "", "the participators of new meeting")
}
