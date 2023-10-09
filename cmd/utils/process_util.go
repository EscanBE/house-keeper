package utils

import (
	"bufio"
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"os"
	"os/exec"
	"sync"
)

func LaunchApp(appName string, args []string, envVars []string, directStd bool) int {
	if directStd {
		return LaunchAppWithDirectStd(appName, args, envVars)
	}
	return LaunchAppWithOutputCallback(appName, args, envVars, nil, nil, nil, nil)
}

func LaunchAppWithOutputCallback(appName string, args []string, envVars []string, stdOutCallback1, stdErrCallBack1, stdOutCallback2, stdErrCallBack2 func(msg string)) int {
	launchCmd := exec.Command(appName, args...)

	launchCmd.Env = envVars
	stdout, errPipeStdout := launchCmd.StdoutPipe()
	if errPipeStdout != nil {
		libutils.PrintfStdErr("problem when getting stdout pipe for %s: %s\n", appName, errPipeStdout.Error())
	}
	stderr, errPipeStderr := launchCmd.StderrPipe()
	if errPipeStderr != nil {
		libutils.PrintfStdErr("problem when getting stderr pipe for %s: %s\n", appName, errPipeStderr.Error())
	}
	stdOutScanner := bufio.NewScanner(stdout)
	stdErrScanner := bufio.NewScanner(stderr)
	err := launchCmd.Start()
	if err != nil {
		libutils.PrintlnStdErr("problem when starting", appName, err)
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
			libutils.PrintlnStdErr("problem when waiting process", appName, err)
			chanEc <- 1
		} else {
			chanEc <- 0
		}
		defer wg.Done()
	}()

	wg.Wait()
	return <-chanEc
}

func LaunchAppWithDirectStd(appName string, args []string, envVars []string) int {
	launchCmd := exec.Command(appName, args...)
	if len(envVars) > 0 {
		launchCmd.Env = envVars
	}
	launchCmd.Stdin = os.Stdin
	launchCmd.Stdout = os.Stdout
	launchCmd.Stderr = os.Stderr
	err := launchCmd.Run()
	if err != nil {
		libutils.PrintfStdErr("problem when running process %s: %s\n", appName, err.Error())
		return 1
	}
	return 0
}
