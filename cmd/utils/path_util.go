package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func HasToolSshPass() bool {
	cmdApp := exec.Command("sshpass", "-V")
	if err := cmdApp.Run(); err != nil {
		return false
	}
	return true
}

func TryReadSymlink(concernPath string) (actualPath string, err error) {
	var file os.FileInfo
	file, err = os.Stat(concernPath)

	if err == nil {
		if file.Mode()&os.ModeSymlink != 0 {
			// not symlink
			actualPath = concernPath
		} else {
			// is symlink
			var resolvedLink string
			resolvedLink, err = os.Readlink(concernPath)
			if err == nil {
				if filepath.IsAbs(resolvedLink) {
					actualPath = resolvedLink
				} else {
					actualPath = path.Join(path.Dir(concernPath), resolvedLink)
				}
				// fmt.Printf("Resolved symlink [%s] => [%s]\n", concernPath, actualPath)
			}
		}
	}

	return
}

func SumDirectorySize(path string, stopAfter int64) (size int64, err error) {
	err = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
			if stopAfter > 0 && size >= stopAfter {
				return errLimitReached
			}
		}
		return err
	})

	return
}

var errLimitReached = fmt.Errorf("limit reached")

func IsErrorLimitSumDirectorySizeReached(err error) bool {
	return err == errLimitReached
}

func HasBinaryName(binaryName string) bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}
