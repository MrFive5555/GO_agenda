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

// clearMeetingCmd represents the clearMeeting command
var clearMeetingCmd = &cobra.Command{
	Use:   "clearMeeting",
	Short: "clear all your meetings",
	Long:  `clear all the meetings that you have created before`,
	Run: func(cmd *cobra.Command, args []string) {

		debugLog("[command] clearMeeting " + strings.Join(args, " "))

		// other argument
		if len(args) > 0 {
			fmt.Println("too many arguments")
			debugLog("too many arguments")
			return
		}

		var state entity.LogState
		entity.GetLogState(&state)

		if state.HasLogin == false {
			fmt.Println("[fail] you haven't logged in any account")
			debugLog("[fail] you haven't logged in any account")
			return
		}

		// 以该用户为 发起者 的会议将被删除
		var meetings entity.MeetingList
		entity.GetMeeting(&meetings)

		flag := false
		toRemove := make([]bool, len(meetings))
		for index, meeting := range meetings {
			if meeting.Sponsors == state.UserName {
				// delete meeting
				toRemove[index] = true
				flag = true
			}
		}

		deleteMeeting(toRemove)

		if flag {
			fmt.Printf("[success] delete all the meetings sponsored by %s\n", state.UserName)
			debugLog("[success] delete all the meetings sponsored by " + state.UserName)
		} else {
			fmt.Printf("[fail] no meeting is sponsored by %s\n", state.UserName)
			debugLog("[fail] no meeting is sponsored by " + state.UserName)
		}

	},
}

func init() {
	rootCmd.AddCommand(clearMeetingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearMeetingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearMeetingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
