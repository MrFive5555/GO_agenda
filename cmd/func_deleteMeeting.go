package cmd

import (
	"fmt"

	"github.com/MrFive5555/GO_agenda/entity"
)

func deleteMeeting(toRemove []bool) {
	var meetings entity.MeetingList
	entity.GetMeeting(&meetings)
	var newMeetings entity.MeetingList
	for i := 0; i < len(meetings); i++ {
		if toRemove[i] {
			fmt.Printf("[success] delete meeting %s\n", meetings[i].Title)
			debugLog("[success] delete meeting %s\n", meetings[i].Title)
		} else {
			newMeetings = append(newMeetings, meetings[i])
		}
	}
	entity.SetMeeting(&newMeetings)
}
