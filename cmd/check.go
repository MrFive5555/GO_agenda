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

	"github.com/MrFive5555/GO_agenda/entity"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check meetings that you participate in",
	Long: `you can check all meetings that you participate in, 
or just check those in a period if specify the start_time and the end_time`,
	Run: func(cmd *cobra.Command, args []string) {

		// 检查参数
		if start == "" {
			start = "0000-01-01-00-00"
		}
		if end == "" {
			end = "9999-12-31-23-59"
		}
		startFormat := timeFormation(start)
		endFormat := timeFormation(end)
		validArgs := true
		infos := []string{"start time", "end time"}
		registerArgs := []string{startFormat, endFormat}
		isvalid := []func(string) bool{isvalidStart, isvalidEnd}
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

		startTime, _ := time.Parse(layout, startFormat)
		endTime, _ := time.Parse(layout, endFormat)
		var meetings entity.MeetingList
		var myMeetings entity.MeetingList
		entity.GetMeeting(&meetings)
		for _, meeting := range meetings {
			mStartFormat := timeFormation(meeting.Start)
			mEndFormat := timeFormation(meeting.End)
			mStart, _ := time.Parse(layout, mStartFormat)
			mEnd, _ := time.Parse(layout, mEndFormat)
			if !(mEnd.Before(startTime) || mStart.After(endTime)) {
				if meeting.Sponsors == me {
					myMeetings = append(myMeetings, meeting)
				}
				mParticipatorsList := strings.Split(meeting.Participators, ",")
				for _, mParticipator := range mParticipatorsList {
					if mParticipator == me {
						myMeetings = append(myMeetings, meeting)
					}
				}
			}
		}

		for key, meeting := range myMeetings {
			fmt.Printf("[meeting %d]\n\ttitle:\t%s\n\tsponsor:\t%s\n\tparticipators:\t%s\n\tstart_time:\t%s\n\tend_time:\t%s\n", key+1, meeting.Title, meeting.Sponsors, meeting.Participators, meeting.Start, meeting.End)
		}
		fmt.Printf("[success] Done! A total of %d meetings have been shown\n", len(myMeetings))
		debugLog("[success] Meetings have been shown")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&start, "start", "s", "", "the start time of new meeting")
	checkCmd.Flags().StringVarP(&end, "end", "e", "", "the end time of new meeting")
}
