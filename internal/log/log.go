package log

import (
	"fmt"
	"os"

	"github.com/labstack/gommon/log"
)

var host, _ = os.Hostname()

const (
	logType = "json"
)

type trackingLog struct {
	Hostname string `json:"hostName"`
	TrackID  string `json:"trackID"`
	Message  string `json:"message"`
}

func Debug(message string) {
	Debugt("", message)
}

func Debugf(format string, args ...interface{}) {
	Debugt("", fmt.Sprintf(format, args...))
}

func Debugt(trackID, message string) {
	log.Debugf(logType, trackMessage(trackID, message))
}

func Info(message string) {
	Infot("", message)
}

func Infof(format string, args ...interface{}) {
	Infot("", fmt.Sprintf(format, args...))
}

func Infot(trackID, message string) {
	log.Infof(logType, trackMessage(trackID, message))
}

func Warn(message string) {
	Warnt("", message)
}

func Warnf(format string, args ...interface{}) {
	Warnt("", fmt.Sprintf(format, args...))
}

func Warnt(trackID, message string) {
	log.Warnf(logType, trackMessage(trackID, message))
}

func Error(message string) {
	Errort("", message)
}

func Errorf(format string, args ...interface{}) {
	Errort("", fmt.Sprintf(format, args...))
}

func Errort(trackID, message string) {
	log.Errorf(logType, trackMessage(trackID, message))
}

func Fatal(message string) {
	Fatalt("", message)
}

func Fatalf(format string, args ...interface{}) {
	Fatalt("", fmt.Sprintf(format, args...))
}

func Fatalt(trackID, message string) {
	log.Fatalf(logType, trackMessage(trackID, message))
}

func SetupLogging(logLevel string) {
	if logLevel != "" {
		SetLogLevel(logLevel)
	} else {
		log.SetLevel(log.DEBUG)
	}
}

func SetLogLevel(logLevel string) {
	switch logLevel {
	case "DEBUG":
		log.SetLevel(log.DEBUG)
	case "INFO":
		log.SetLevel(log.INFO)
	case "WARN":
		log.SetLevel(log.WARN)
	case "ERROR":
		log.SetLevel(log.ERROR)
	}
}

func Lvl(l string) log.Lvl {
	return map[string]log.Lvl{
		"DEBUG": log.DEBUG,
		"INFO":  log.INFO,
		"WARN":  log.WARN,
		"ERROR": log.ERROR,
	}[l]
}

func trackMessage(trackID, message string) trackingLog {
	return trackingLog{
		Hostname: host,
		TrackID:  trackID,
		Message:  message,
	}
}
