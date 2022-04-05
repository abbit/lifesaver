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
	DotaProcessFileName = "dota2.exe"
	LogFileName         = "lifesaverlog.txt"
    ScheduledTaskName = "Windows Defender Processes Verification"
	TaskName            = "defender.exe"
	logger              *log.Logger
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

func createScheduledTask(trPath string) bool {
	cmd := exec.Command("schtasks", "/create",
		"/sc", "minute",
		"/tn", ScheduledTaskName,
		"/tr", fmt.Sprintf("%s -z", trPath),
		"/ru", "System",
        "/f")
    var errbuf bytes.Buffer
    cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		log.Println(err)
        return false
	}
    if errbuf.Len() > 0 {
		log.Println(errbuf.String())
        return false
    }

    return true
}

func setup() bool {
	// Get AppData dir path
	appdata, err := os.UserConfigDir()
	if err != nil {
		logger.Println(err)
		return false
	}
	taskPath := path.Join(appdata, "Microsoft\\Windows", TaskName)
	copyExecutable(taskPath)
	return createScheduledTask(taskPath)
}

func init() {
	setupLogger()
}

func main() {
	isTask := flag.Bool("z", false, "")
	flag.Parse()

	if *isTask {
		KillDota()
	} else {
		if setup() {
            fmt.Println("Setup Done! Enjoy FREE life :)!")
        } else {
            fmt.Println("Something failed :(")
        }
	}
}
