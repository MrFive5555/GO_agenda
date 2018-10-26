// load the status of agenda
package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	CACHE_DIR    = "./.cache/"
	USER_JSON    = CACHE_DIR + "User.json"
	MEETING_JSON = CACHE_DIR + "Meeting.json"
	LOG_JSON     = CACHE_DIR + "Log.json"
)

type User struct {
	UserName  string `json:"name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

// type Time struct {
// 	Year  int `json:"Year"`
// 	Mouth int `json:"Mouth"`
// 	Day   int `json:"Day"`
// 	Hour  int `json:"Hour"`
// 	Min   int `json:"Min"`
// }

type Meeting struct {
	Title         string `json:"Title"`
	Sponsors      string `json:"Sponsors"`
	Participators string `json:"Participators"`
	Start         string `json:"Start"`
	End           string `json:"End"`
}

type UserList []User
type MeetingList []Meeting

type LogState struct {
	HasLogin bool   `json:"HasLogin"`
	UserName string `json:"UserName"`
}

var username, password, email, telephone string
var title, sponsor, participators, start, end string

const layout string = "2006-01-02 15:04:05"

func GetUsers(usersPtr *UserList) error {
	return loadJSON(USER_JSON, usersPtr)
}
func GetMeeting(meetingsPtr *MeetingList) error {
	return loadJSON(MEETING_JSON, meetingsPtr)
}
func GetLogState(statePtr *LogState) error {
	return loadJSON(LOG_JSON, statePtr)
}
func SetUsers(usersPtr *UserList) error {
	return saveJSON(USER_JSON, usersPtr)
}
func SetMeeting(meetingsPtr *MeetingList) error {
	return saveJSON(MEETING_JSON, meetingsPtr)
}
func SetLogState(statePtr *LogState) error {
	return saveJSON(LOG_JSON, statePtr)
}

func loadJSON(file string, dataPtr interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(data), dataPtr)
	if err != nil {
		return err
	}
	return nil
}
func saveJSON(file string, dataPtr interface{}) error {
	b, err := json.Marshal(dataPtr)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, b, 0664)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	os.Mkdir(CACHE_DIR, 0775)
	checkCache := func(name string, default_context string) {
		_, err := os.Stat(name)
		if err != nil {
			ioutil.WriteFile(name, []byte(default_context), 0664)
		}
	}
	checkCache(USER_JSON, "{}")
	checkCache(MEETING_JSON, "{}")
	checkCache(LOG_JSON, `{"HasLogin":false,"UserName":""}`)
}
