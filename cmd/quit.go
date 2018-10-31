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

	"github.com/spf13/cobra"
)

// quitCmd represents the quit command
var quitCmd = &cobra.Command{
	Use:   "quit",
	Short: "quit a meeting",
	Long:  `quit a meeting by specifying its title with -t`,
	Run: func(cmd *cobra.Command, args []string) {

		debugLog("[command] quit " + strings.Join(args, " "))

		// other argument
		if len(args) > 0 {
			fmt.Println("too many arguments")
			debugLog("too many arguments")
			return
		}

		var state LogState
		GetLogState(&state)

		if state.HasLogin == false {
			fmt.Println("[fail] you haven't logged in any account")
			debugLog("[fail] you haven't logged in any account")
			return
		}

		// 以该用户为 参与者 的会议将从 参与者 列表中移除该用户
		// 若因此造成会议 参与者 人数为0，则会议也将被删除。
		var meetings MeetingList
		GetMeeting(&meetings)

		toRemove := make([]bool, len(meetings))
		for index, meeting := range meetings {
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
						SetMeeting(&meetings)
						fmt.Printf("[success] delete participator %s from meeting %s\n", state.UserName, meeting.Title)
						debugLog("[success] delete participator " + state.UserName + " from meeting " + meeting.Title)
					}
				}
			}
		}

		deleteMeeting(toRemove)

		fmt.Println("[success] quit meeting successfully")
		debugLog("[success] quit meeting successfully")

	},
}

func init() {
	rootCmd.AddCommand(quitCmd)

	quitCmd.Flags().StringVarP(&quitTitle, "title", "t", "", "the title of meeting to quit")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
