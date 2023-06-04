package utils

import (
	"os/exec"
)

func HasToolSshPass() bool {
	cmdApp := exec.Command("sshpass", "-V")
	if err := cmdApp.Run(); err != nil {
		return false
	}
	return true
}
