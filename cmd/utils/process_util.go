package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func LaunchApp(appName string, args []string, envVars []string) int {
	return LaunchAppWithOutputCallback(appName, args, envVars, nil, nil)
}

func LaunchAppWithOutputCallback(appName string, args []string, envVars []string, stdOutCallback, stdErrCallBack func(msg string)) int {
	rsyncCmd := exec.Command(appName, args...)

	rsyncCmd.Env = envVars
	stdout, _ := rsyncCmd.StdoutPipe()
	stderr, _ := rsyncCmd.StderrPipe()
	rsyncStdOutScanner := bufio.NewScanner(stdout)
	rsyncStdErrScanner := bufio.NewScanner(stderr)
	err := rsyncCmd.Start()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "problem when starting", appName, err)
	}

	var chanEc = make(chan int, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			oScan := rsyncStdOutScanner.Scan()
			eScan := rsyncStdErrScanner.Scan()
			if oScan {
				msg := rsyncStdOutScanner.Text()
				fmt.Println(msg)
				if stdOutCallback != nil {
					stdOutCallback(msg)
				}
			}
			if eScan {
				msg := rsyncStdErrScanner.Text()
				_, _ = fmt.Fprintln(os.Stderr, msg)
				if stdErrCallBack != nil {
					stdErrCallBack(msg)
				}
			}
			if !oScan && !eScan {
				break
			}
		}
		err = rsyncCmd.Wait()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "problem when waiting process", appName, err)
			chanEc <- 1
		} else {
			chanEc <- 0
		}
		defer wg.Done()
	}()

	wg.Wait()
	return <-chanEc
}
