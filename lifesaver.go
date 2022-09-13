package main

import (
	"bytes"
	"flag"
	"fmt"
	winapi "github.com/winlabs/gowin32"
	"log"
	"os"
	"os/exec"
	"path"
)

var (
	LogFileName       = "lifesaverlog.txt"
	ScheduledTaskName = "Windows Defender Processes Verification"
	TaskPath          = path.Join("C:", "Windows", "defender.exe")
	DefaultBannedExes = "dota2.exe"
	logger            *log.Logger
)

func initLogger() {
	logFilePath := LogFileName
	logFilePath = path.Join("C:", LogFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		logger = log.Default()
		logger.Println("Error while opening log file", err)
		return
	}

	logger = log.New(logFile, "", log.LstdFlags)
}

func killProcesses(processes []string) {
	processMap := make(map[string]bool)
	for _, process := range processes {
		processMap[process] = true
	}

	procs, err := winapi.GetProcesses()
	if err != nil {
		logger.Println("Error while getting list of processes: ", err)
		return
	}

	for _, proc := range procs {
		if processMap[proc.ExeFile] {
			err = winapi.KillProcess(proc.ProcessID, 322)
			if err != nil {
				logger.Printf("Can't kill %s: %s\n", proc.ExeFile, err)
			}
			logger.Printf("%s successfully killed, GG!\n", proc.ExeFile)
		}
	}
}

func copyExecutable(toPath string) {
	// Get path of this executable
	exe, err := os.Executable()
	if err != nil {
		logger.Println(err)
		return
	}

	// If already there, then no need to copy
	if toPath == exe {
		return
	}

	// Read this executable
	r, err := os.Open(exe)
	if err != nil {
		logger.Println(err)
		return
	}
	defer r.Close()

	// copy executable to dest dir
	w, err := os.Create(toPath)
	if err != nil {
		logger.Println(err)
		return
	}
	defer w.Close()
	w.ReadFrom(r)
}

// TODO: create a way to add apps to banned list (append them to /tr option)

func createScheduledTask(trPath string) bool {
	cmd := exec.Command("schtasks", "/create",
		"/sc", "minute",
		"/tn", ScheduledTaskName,
		"/tr", fmt.Sprintf("%s -z %s", trPath, DefaultBannedExes),
		"/ru", "System",
		"/f")
	var errbuf bytes.Buffer
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		logger.Println(err)
		return false
	}
	if errbuf.Len() > 0 {
		logger.Println(errbuf.String())
		return false
	}

	return true
}

func setup(dstPath string) {
	copyExecutable(dstPath)
	if createScheduledTask(dstPath) {
		logger.Println("Setup Done! Enjoy FREE life :)!")
	} else {
		logger.Println("Something failed :(")
	}
}

func main() {
	initLogger()

	isTask := flag.Bool("z", false, "")
	flag.Parse()

	if *isTask {
		killProcesses(flag.Args())
	} else {
		setup(TaskPath)
	}
}
