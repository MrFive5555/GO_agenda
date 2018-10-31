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

// createMeetingCmd represents the createMeeting command
var createMeetingCmd = &cobra.Command{
	Use:   "createMeeting",
	Short: "create a meeting",
	Long: `create a new meeting to the agenda, the usage of it
	you should specify the title, participators(separated by commas, no space after commas), start time and end time`,
	Run: func(cmd *cobra.Command, args []string) {

		debugLog("[command] createMeeting -t " + title + " -p " + participators + " -s " + start + " -e " + end + " " + strings.Join(args, " "))

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

		// make the intput to target format
		startFormat := timeFormation(start)
		endFormat := timeFormation(end)
		if startFormat == "" || endFormat == "" {
			return
		}

		// 检查非法参数
		validArgs := true
		infos := []string{"title", "participators", "start time", "end time"}
		registerArgs := []string{title, participators, startFormat, endFormat}
		isvalid := []func(string) bool{isvalidTitle, isvalidParticipators, isvalidStart, isvalidEnd}
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

		// split participators
		participatorsList := strings.Split(participators, ",")

		// parse time to date
		startTime, _ := time.Parse(layout, startFormat)
		endTime, _ := time.Parse(layout, endFormat)

		if startTime.After(endTime) || startTime.Equal(endTime) {
			fmt.Println("[fail] start time must before end time")
			debugLog("[fail] start time must before end time")
			return
		}

		// 检查是否重名
		var meetings MeetingList
		GetMeeting(&meetings)
		for _, meeting := range meetings {
			if meeting.Title == title {
				fmt.Printf("[fail] there has been a meeting with the same title %s\n", title)
				debugLog("[fail] there has been a meeting with the same title " + title)
				return
			}
		}

		// 不允许包含未注册用户
		var users UserList
		GetUsers(&users)
		for _, participator := range participatorsList {
			exist := false
			for _, user := range users {
				if user.UserName == participator {
					exist = true
					break
				}
			}
			if !exist {
				fmt.Printf("[fail] participator %s does not exist\n", participator)
				debugLog("[fail] participator " + participator + " does not exist")
				return
			}
		}

		// 任何用户都无法分身参加多个会议。
		// 如果用户已有的会议安排（作为发起者或参与者）
		// 与将要创建的会议在时间上重叠 （允许仅有端点重叠的情况），则无法创建该会议
		for _, participator := range participatorsList {
			for _, meeting := range meetings {
				mStartFormat := timeFormation(meeting.Start)
				mEndFormat := timeFormation(meeting.End)
				mStart, _ := time.Parse(layout, mStartFormat)
				mEnd, _ := time.Parse(layout, mEndFormat)
				if !((mEnd.Before(startTime) || mEnd.Equal(startTime)) || (mStart.After(endTime) || mStart.Equal(endTime))) {
					if meeting.Sponsors == participator {
						fmt.Printf("[fail] participator %s is a busy sponsor in another meeting (%s)\n", participator, meeting.Title)
						debugLog("[fail] participator " + participator + " is a busy sponsor in another meeting (" + meeting.Title + ")")
						return
					}
					mParticipatorsList := strings.Split(meeting.Participators, ",")
					for _, mParticipator := range mParticipatorsList {
						if mParticipator == participator {
							fmt.Printf("[fail] participator %s is a busy participator in another meeting (%s)\n", participator, meeting.Title)
							debugLog("[fail] participator " + participator + " is a busy participator in another meeting (" + meeting.Title + ")")
							return
						}
					}
				}
			}
		}

		meetings = append(meetings, Meeting{
			title,
			state.UserName,
			participators,
			start,
			end,
		})
		SetMeeting(&meetings)
		fmt.Printf("[success] new meeting %s has been added\n", title)
		debugLog("[success] new meeting " + title + " has been added")

	},
}

func init() {
	rootCmd.AddCommand(createMeetingCmd)

	createMeetingCmd.Flags().StringVarP(&title, "title", "t", "", "the title of new meeting")
	createMeetingCmd.Flags().StringVarP(&participators, "participators", "p", "", "the participators of new meeting")
	createMeetingCmd.Flags().StringVarP(&start, "start", "s", "", "the start time of new meeting")
	createMeetingCmd.Flags().StringVarP(&end, "end", "e", "", "the end time of new meeting")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createMeetingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createMeetingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func isvalidTitle(title string) bool {
	return title != ""
}

func isvalidParticipators(participators string) bool {
	// should separated by commas without space
	for i := 0; i < len(participators); i++ {
		if participators[i] == ' ' {
			fmt.Printf("[fail] %s contains space, no space is expected\n", participators)
			debugLog("[fail] " + participators + " contains space, no space is expected\n")
			return false
		}
	}
	return participators != ""
}

func isvalidStart(start string) bool {
	// handle time format error
	_, err := time.Parse(layout, start)
	if err != nil {
		// fmt.Printf("[fail] invalid format of start time %s, should be like 2006-01-02-15-04\n", start)
		fmt.Printf("[fail] %s\n", err.Error())
		debugLog("[fail] " + err.Error())
		return false
	}
	return true
}

func isvalidEnd(end string) bool {
	// handle time format error
	_, err := time.Parse(layout, end)
	if err != nil {
		// fmt.Printf("[fail] invalid format of end time %s, should be like 2006-01-02-15-04\n", end)
		fmt.Printf("[fail] %s\n", err.Error())
		debugLog("[fail] " + err.Error())
		return false
	}
	return true
}

func timeFormation(origin string) string {
	test := strings.Split(origin, "-")
	if len(test) != 5 {
		fmt.Printf("[fail] invalid format of end time %s, should be like 2006-01-02-15-04\n", end)
		debugLog("[fail] invalid format of end time " + end + ", should be like 2006-01-02-15-04")
		return ""
	}
	return test[0] + "-" + test[1] + "-" + test[2] + " " + test[3] + ":" + test[4] + ":00"
}

// func myTimeParse(origin string) ([]int, string) {
// 	var arr []int
// 	str := ""
// 	test := strings.Split(origin, "-")
// 	if len(test) != 5 {
// 		return arr, ""
// 	}
// 	year, err := strconv.Atoi(test[0])
// 	if err != nil {
// 		str = "error while parsing year\n"
// 		return arr, str
// 	}
// 	arr[0] = year

// 	month, err := strconv.Atoi(test[1])
// 	if err != nil {
// 		str = "error while parsing month\n"
// 		return arr, str
// 	}
// 	arr[1] = month

// 	day, err := strconv.Atoi(test[2])
// 	if err != nil {
// 		str = "error while parsing day\n"
// 		return arr, str
// 	}
// 	arr[2] = day

// 	hour, err := strconv.Atoi(test[3])
// 	if err != nil {
// 		str = "error while parsing hour\n"
// 		return arr, str
// 	}
// 	arr[3] = hour

// 	minute, err := strconv.Atoi(test[4])
// 	if err != nil {
// 		str = "error while parsing minute\n"
// 		return arr, str
// 	}
// 	arr[4] = minute
// 	arr[5] = 0
// 	return arr, str
// }
