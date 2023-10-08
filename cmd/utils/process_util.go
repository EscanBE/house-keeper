package utils

import (
	"bufio"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"os"
	"os/exec"
	"sync"
)

func LaunchApp(appName string, args []string, envVars []string) int {
	return LaunchAppWithOutputCallback(appName, args, envVars, nil, nil, nil, nil)
}

func LaunchAppWithOutputCallback(appName string, args []string, envVars []string, stdOutCallback1, stdErrCallBack1, stdOutCallback2, stdErrCallBack2 func(msg string)) int {
	launchCmd := exec.Command(appName, args...)

	launchCmd.Env = envVars
	stdout, _ := launchCmd.StdoutPipe()
	stderr, _ := launchCmd.StderrPipe()
	stdOutScanner := bufio.NewScanner(stdout)
	stdErrScanner := bufio.NewScanner(stderr)
	err := launchCmd.Start()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "problem when starting", appName, err)
	}

	var chanEc = make(chan int, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			oScan := stdOutScanner.Scan()
			eScan := stdErrScanner.Scan()
			if oScan {
				msg := stdOutScanner.Text()
				fmt.Println(msg)
				if stdOutCallback1 != nil {
					stdOutCallback1(msg)
				}
				if stdOutCallback2 != nil {
					stdOutCallback2(msg)
				}
			}
			if eScan {
				msg := stdErrScanner.Text()
				libutils.PrintlnStdErr(msg)
				if stdErrCallBack1 != nil {
					stdErrCallBack1(msg)
				}
				if stdErrCallBack2 != nil {
					stdErrCallBack2(msg)
				}
			}
			if !oScan && !eScan {
				break
			}
		}
		err = launchCmd.Wait()
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
