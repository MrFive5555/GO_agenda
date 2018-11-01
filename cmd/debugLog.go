package cmd

import (
	"fmt"
	"log"
	"os"
)

var debugLog func(format string, arg ...interface{})

func init() {
	debugLog = (func() func(string, ...interface{}) {
		logfile, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(3)
		}
		debugLogger := log.New(logfile, "[debug]", log.LstdFlags)

		return func(format string, arg ...interface{}) {
			str := fmt.Sprintf(format, arg...)
			debugLogger.Println(str)
		}
	})()
}
