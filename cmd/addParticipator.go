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
	"time"

	"github.com/spf13/cobra"
)

// addParticipatorCmd represents the addParticipator command
var addParticipatorCmd = &cobra.Command{
	Use:   "addParticipator",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("addParticipator called")

		// 检查非法参数
		validArgs := true
		infos := []string{"title", "participators"}
		registerArgs := []string{title, participators}
		isvalid := []func(string) bool{isvalidTitle, isvalidParticipators}
		for i, info := range infos {
			if !isvalid[i](registerArgs[i]) {
				validArgs = false
				fmt.Printf("[fail] the Field %s is invalid\n", info)
			}
		}
		if !validArgs {
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
			return 
		}

		// 检查自己发起的主题为title的会议是否存在，并获取该会议信息
		var meetings MeetingList
		GetMeeting(&meetings)
		var myMeeting Meeting
		validMeeting := false
		for _, m := range meetings {
			if (m.Title == title) && (m.Sponsors == me) {
				validMeeting = true
				myMeeting = m
				break
			}
		}
		if !validMeeting {
			fmt.Printf("[fail] %s does not have sponsored the meeting %s\n", me, title)
			return
		}

		// 检查用户是否已注册
		var users UserList
		GetUsers(&users)
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
				return
			}
		}

		// 任何用户都无法分身参加多个会议。
		// 如果用户已有的会议安排（作为发起者或参与者）
		// 与将要参加的会议在时间上重叠 （允许仅有端点重叠的情况），则无法参加该会议
		startFormat := timeFormation(myMeeting.Start)
		endFormat := timeFormation(myMeeting.End)
		startTime, _ := time.Parse(layout, startFormat)
		endTime, _ := time.Parse(layout, endFormat)

		for _, participator := range participatorsList {
			for _, meeting := range meetings {
				mStartFormat := timeFormation(meeting.Start)
				mEndFormat := timeFormation(meeting.End)
				mStart, _ := time.Parse(layout, mStartFormat)
				mEnd, _ := time.Parse(layout, mEndFormat)
				if !((mEnd.Before(startTime) || mEnd.Equal(startTime)) || (mStart.After(endTime) || mStart.Equal(endTime))) {
					if meeting.Sponsors == participator {
						fmt.Printf("[fail] participator %s is a busy sponsor in another meeting (%s)\n", participator, meeting.Title)
						return
					}
					mParticipatorsList := strings.Split(meeting.Participators, ",")
					for _, mParticipator := range mParticipatorsList {
						if mParticipator == participator {
							fmt.Printf("[fail] participator %s is a busy participator in another meeting (%s)\n", participator, meeting.Title)
							return
						}
					}
				}
			}
		}

		// 如无错误，增加会议参与者
		for i, m := range meetings {
			if m.Title == title {
				meetings[i].Participators += "," 
				meetings[i].Participators += participators
				break
			}
		}
		SetMeeting(&meetings)
		fmt.Printf("[success] new participator(s) %s has(have) been added into the %s\n", participators, title)

	},
}

func init() {
	rootCmd.AddCommand(addParticipatorCmd)

	addParticipatorCmd.Flags().StringVarP(&title, "title", "t", "", "the title of new meeting")
	addParticipatorCmd.Flags().StringVarP(&participators, "participators", "p", "", "the participators of new meeting")
}
