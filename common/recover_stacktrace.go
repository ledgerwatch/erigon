package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/ledgerwatch/erigon/log"
)

var sigc chan os.Signal

func GetSigC(sig *chan os.Signal) {
	sigc = *sig
}

func prettyTime() string {
	time := fmt.Sprintf("%v", time.Now())
	return strings.Replace(time[:19], " ", "-", 1)
}

func CheckForCrashes() {
	ex, err := os.Executable()
	if err != nil {
		log.Error(err.Error())
		return
	}
	binPath := filepath.Dir(ex)
	crashReportDir := filepath.Join(binPath[:len(binPath)-10], "crashreports")
	f, err := os.Open(crashReportDir)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fileInfo, err := f.ReadDir(-1)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, v := range fileInfo {
		if !v.IsDir() {
			msg := fmt.Sprintf("Crashes From Previous Boots Detected. Find the stack trace in %v",
				crashReportDir)
			log.Warn(msg)
			return
		}
	}
}

func RecoverStackTrace(r interface{}) {
	if r != nil {
		stack := string(debug.Stack())
		WriteStackTraceOnPanic(stack)
		sigc <- syscall.SIGINT
	}
}

func WriteStackTraceOnPanic(stack string) {
	ex, err := os.Executable()
	if err != nil {
		log.Error(err.Error())
		return
	}
	binPath := filepath.Dir(ex)
	fileName := filepath.Join(binPath[:len(binPath)-10], "crashreports", prettyTime()+".txt")
	f, errFs := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if errFs != nil {
		log.Error(errFs.Error())
		f.Close()
		return
	}
	f.WriteString(stack)
	f.Sync()
	f.Close()
}
