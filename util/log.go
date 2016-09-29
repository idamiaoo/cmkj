package util

import (
	"os"

	"github.com/op/go-logging"
)

var consoleFormat = logging.MustStringFormatter(
	`%{color}%{time:15:04:05} %{shortfile} %{shortfunc} ▶ %{level} %{color:reset} %{message}`,
)

var fileFormat = logging.MustStringFormatter(
	`%{time:2006-01-02 15:04:05} %{shortfile} %{shortfunc} ▶ %{level} %{message}`,
)

/*
var allFileFormat = logging.MustStringFormatter(
	`[%{module}] %{time:2006-01-02 15:04:05} %{shortfile} %{shortfunc} ▶ %{level} %{message}`,
)
*/

var Log = logging.MustGetLogger("log")

func InitLog(logname, loglevel string) {
	level, err := logging.LogLevel(loglevel)
	if err != nil {
		panic(err)
	}
	/*
		allFile, err := os.OpenFile(logname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	*/

	logFile, err := os.OpenFile(logname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
	//allFileBackend := logging.NewLogBackend(allFile, "", 0)
	logFileBackend := logging.NewLogBackend(logFile, "", 0)

	consoleFmt := logging.NewBackendFormatter(consoleBackend, consoleFormat)
	//allFileFmt := logging.NewBackendFormatter(allFileBackend, allFileFormat)
	logFileFmt := logging.NewBackendFormatter(logFileBackend, fileFormat)

	consoleLevel := logging.AddModuleLevel(consoleBackend)
	//allFileLevel := logging.AddModuleLevel(allFileBackend)
	logFileLevel := logging.AddModuleLevel(logFileBackend)

	consoleLevel.SetLevel(level, "")
	//allFileLevel.SetLevel(level, "")
	logFileLevel.SetLevel(level, "")

	logBackend := logging.SetBackend(consoleFmt, logFileFmt)

	Log.SetBackend(logBackend)
}
