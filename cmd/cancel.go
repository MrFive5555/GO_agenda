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

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "cancel a meeting of your own",
	Long: `specify the title and you will cancel the meeting you sponsor`,

	Run: func(cmd *cobra.Command, args []string) {

		// 检查非法参数
		if !isvalidTitle(title) {
			fmt.Printf("[fail] the Field title is invalid\n")
			debugLog("[fail] the Field title is invalid")
			return
		}

		// 获得当前登录用户信息
		var state LogState
		var me string
		GetLogState(&state)
		if state.HasLogin {
			me = state.UserName
		} else {
			fmt.Printf("[fail] you haven't loged in\n")
			debugLog("[fail] you haven't logged in")
			return 
		}

		// 检查自己发起的主题为title的会议是否存在，并获取该会议信息
		var meetings MeetingList
		GetMeeting(&meetings)
		toRemove := make([]bool, len(meetings))
		validMeeting := false
		for i, m := range meetings {
			if (m.Title == title) && (m.Sponsors == me) {
				validMeeting = true
				toRemove[i] = true
				break
			}
		}
		if !validMeeting {
			fmt.Printf("[fail] %s does not have sponsored the meeting %s\n", me, title)
			debugLog("[fail] " + me + " does not have sponsored the meeting " + title)
			return
		}
		deleteMeeting(toRemove)
		
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	cancelCmd.Flags().StringVarP(&title, "title", "t", "", "the title of new meeting")
}
